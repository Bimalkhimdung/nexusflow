variable "environment" {
  type = string
}

resource "aws_s3_bucket" "attachments" {
  bucket = "nexusflow-${var.environment}-attachments-${random_id.suffix.hex}"

  tags = {
    Name = "nexusflow-${var.environment}-attachments"
  }
}

resource "random_id" "suffix" {
  byte_length = 4
}

resource "aws_s3_bucket_ownership_controls" "attachments" {
  bucket = aws_s3_bucket.attachments.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "attachments" {
  depends_on = [aws_s3_bucket_ownership_controls.attachments]

  bucket = aws_s3_bucket.attachments.id
  acl    = "private"
}

output "bucket_name" {
  value = aws_s3_bucket.attachments.id
}
