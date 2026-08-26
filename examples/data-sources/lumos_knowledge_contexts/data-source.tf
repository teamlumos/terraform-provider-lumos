data "lumos_knowledge_contexts" "my_knowledgecontexts" {
  page_size  = 10
  page_token = "...my_page_token..."
  product_scopes = [
    "ALBUS_DEEP_RESEARCH"
  ]
  title_search = "...my_title_search..."
}