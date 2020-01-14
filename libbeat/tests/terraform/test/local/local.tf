variable "install_dir" {
  type = string
}

provider "local" {
  version = "~> 1.4"
}

data "local_file" "source" {
  filename = "file.txt"
}

resource "local_file" "file_txt" {
  content         = "${data.local_file.source.content}"
  filename        = "${var.install_dir}/file.txt"
  file_permission = "0644"
}

resource "local_file" "empty_file" {
  filename        = "${var.install_dir}/empty_file"
  file_permission = "0644"
}

output "file_path" {
  value = "${local_file.file_txt.filename}"
}

output "empty_file_path" {
  value = "${local_file.empty_file.filename}"
}
