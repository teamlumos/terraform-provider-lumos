#!/usr/bin/env perl
# Makes `terraform show` output comparable across runs and environments:
# replaces UUIDs with <uuid-N>, numbered in order of first appearance within
# a case, the run's name suffix with <run>, and the HTTP version in API errors
# with <version>, and puts the elements of set attributes in a stable order.
#
#   normalize.pl <uuid map file> <run id> [<set attributes file>] < plan.txt > normalized.txt
#
# The map file carries the numbering from one snapshot of a case to the next,
# so the same object keeps the same number in every step. The set attributes
# file names the provider's set attributes, one per line.
use strict;
use warnings;
use sort 'stable';

my ($map_file, $run_id, $sets_file) = @ARGV;
binmode STDIN, ':encoding(UTF-8)';
binmode STDOUT, ':encoding(UTF-8)';

my $uuid = qr/\b[0-9a-fA-F]{8}-(?:[0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}\b/;

my %ids;
if (open my $in, '<', $map_file) {
    while (<$in>) {
        chomp;
        my ($id, $n) = split / /;
        $ids{$id} = $n;
    }
    close $in;
}
my $next = 1 + keys %ids;

my %set_attributes;
if (defined $sets_file && open my $in, '<', $sets_file) {
    while (<$in>) {
        chomp;
        $set_attributes{$_} = 1 if length;
    }
    close $in;
}

my $text = do { local $/; <STDIN> };

# API errors dump the HTTP request and response, with CRLF line endings. Their
# headers (host, date, user agent, rate-limit counters) differ between runs, so
# keep only the request line or status line and the body. Terraform wraps long
# lines, so the request line or a header can take more than one line; the
# headers end at the first blank line. The HTTP version depends on what's in
# front of the API: a local API answers HTTP/1.1, and production HTTP/2.0.
$text =~ s/\r//g;
$text =~ s{^(\*\*Request\*\*:\n(?:.*\n)*?.*HTTP/[\d.]+\n)(?:.+\n)*}{$1}mg;
$text =~ s{^(\*\*Response\*\*:\nHTTP/[\d.]+ \d{3}.*\n)(?:.+\n)*}{$1}mg;
$text =~ s{HTTP/(?:[\d.]+|<version>)}{HTTP/<version>}g;

# Terraform prints a set's elements sorted by value, and the values hold UUIDs
# that are new on every run, so the same plan can list them in a different
# order each time. Sort them by their text with the UUIDs replaced by the
# numbers they already have, before numbering the new ones.
$text = join '', sort_sets(split /^/m, $text) if %set_attributes;

$text =~ s{($uuid)}{'<uuid-' . ($ids{lc $1} //= $next++) . '>'}ge;
$text =~ s/\Q$run_id\E/<run>/gi;
$text =~ s/[ \t]+$//mg;

print $text;

open my $out, '>', $map_file or die "can't write $map_file: $!";
print {$out} "$_ $ids{$_}\n" for sort { $ids{$a} <=> $ids{$b} } keys %ids;
close $out;

# Sorts the elements of every multi-line set attribute in @lines, including
# sets nested in another set's elements. A set opens with `name = [` and
# closes with `]` in the name's column; each element starts two columns in,
# at an action marker, or four columns in without one.
sub sort_sets {
    my @lines = @_;
    my @out;
    for (my $i = 0; $i < @lines; $i++) {
        my ($indent, $action, $name) = $lines[$i] =~ /^( *)((?:[-+~] )?)(\w+) *= \[\n$/;
        my $col = defined $name ? length($indent) + length($action) : 0;
        my $end = defined $name && $set_attributes{$name}
            ? first_index($i + 1, qr/^ {$col}\]/, @lines) : -1;
        if ($end < 0) {
            push @out, $lines[$i];
            next;
        }
        my @elements;
        for my $line (@lines[$i + 1 .. $end - 1]) {
            my $starts = $line =~ /^ {@{[$col + 2]}}[-+~] / || $line =~ /^ {@{[$col + 4]}}[^ }\]]/;
            if ($starts || !@elements) {
                push @elements, [$line];
            } else {
                push @{$elements[-1]}, $line;
            }
        }
        @elements = map { [sort_sets(@$_)] } @elements;
        push @out, $lines[$i], (map { @$_ } sort { sort_key($a) cmp sort_key($b) } @elements), $lines[$end];
        $i = $end;
    }
    return @out;
}

sub first_index {
    my ($from, $pattern, @lines) = @_;
    for my $i ($from .. $#lines) {
        return $i if $lines[$i] =~ $pattern;
    }
    return -1;
}

# An element's text with each UUID replaced by its number, zero-padded so
# numbers sort in order, or by ~ when it has none yet. Notes such as
# "# (1 unchanged element hidden)" sort last.
sub sort_key {
    my ($element) = @_;
    my $key = join '', @$element;
    $key =~ s{($uuid)}{exists $ids{lc $1} ? sprintf('%06d', $ids{lc $1}) : '~'}ge;
    $key =~ s{<uuid-(\d+)>}{sprintf('%06d', $1)}ge;
    return ($key =~ /^\s*#/ ? '1' : '0') . $key;
}
