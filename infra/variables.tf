variable "location" {
  type        = string
  default     = "westeurope"
}

variable "project_name" {
  type        = string
  default     = "devops-project-sk"
}

variable "acr_sku" {
  type = string
  default = "Basic"
}

variable "vm_size" {
  type = string
  default = "Standard_B1s"
}
