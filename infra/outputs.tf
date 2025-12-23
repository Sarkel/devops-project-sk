output "resource_group_name" {
  value = azurerm_resource_group.rg.name
}

output "acr_login_server" {
  value = azurerm_container_registry.acr.login_server
}

output "acr_admin_username" {
  value = azurerm_container_registry.acr.admin_username
}

output "vm_public_ip" {
  value = azurerm_public_ip.pip.ip_address
}

output "vm_private_key" {
  value = tls_private_key.ssh.private_key_openssh
  sensitive = true
}