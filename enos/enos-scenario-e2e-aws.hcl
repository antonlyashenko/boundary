# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This scenario requires access to the boundary team's test AWS account
scenario "e2e_aws" {
  terraform_cli = terraform_cli.default
  terraform     = terraform.default
  providers = [
    provider.aws.default,
    provider.enos.default
  ]

  matrix {
    builder = ["local", "crt"]
  }

  locals {
    aws_ssh_private_key_path = abspath(var.aws_ssh_private_key_path)
    boundary_install_dir     = abspath(var.boundary_install_dir)
    local_boundary_dir       = abspath(var.local_boundary_dir)
    build_path = {
      "local" = "/tmp",
      "crt"   = var.crt_bundle_path == null ? null : abspath(var.crt_bundle_path)
    }
    tags = merge({
      "Project Name" : var.project_name
      "Project" : "Enos",
      "Environment" : "ci"
    }, var.tags)
  }

  step "find_azs" {
    module = module.az_finder

    variables {
      instance_type = [
        var.worker_instance_type,
        var.controller_instance_type
      ]
    }
  }

  step "create_db_password" {
    module = module.random_stringifier
  }

  step "build_boundary" {
    module = matrix.builder == "crt" ? module.build_crt : module.build_local

    variables {
      path = local.build_path[matrix.builder]
    }
  }

  step "create_base_infra" {
    module = module.infra

    variables {
      availability_zones = step.find_azs.availability_zones
      common_tags        = local.tags
    }
  }

  step "create_boundary_cluster" {
    module = module.boundary
    depends_on = [
      step.create_base_infra,
      step.build_boundary
    ]

    variables {
      boundary_install_dir     = local.boundary_install_dir
      common_tags              = local.tags
      controller_instance_type = var.controller_instance_type
      controller_count         = var.controller_count
      db_pass                  = step.create_db_password.string
      kms_key_arn              = step.create_base_infra.kms_key_arn
      local_artifact_path      = step.build_boundary.artifact_path
      ubuntu_ami_id            = step.create_base_infra.ami_ids["ubuntu"]["amd64"]
      vpc_id                   = step.create_base_infra.vpc_id
      worker_count             = var.worker_count
      worker_instance_type     = var.worker_instance_type
    }
  }

  step "create_tag1" {
    module = module.random_stringifier
  }

  step "create_tag1_inputs" {
    module     = module.generate_aws_host_tag_vars
    depends_on = [step.create_tag1]

    variables {
      tag_name  = step.create_tag1.string
      tag_value = "true"
    }
  }

  step "create_targets_with_tag1" {
    module     = module.target
    depends_on = [step.create_base_infra]

    variables {
      ami_id               = step.create_base_infra.ami_ids["ubuntu"]["amd64"]
      aws_ssh_keypair_name = var.aws_ssh_keypair_name
      enos_user            = var.enos_user
      instance_type        = var.target_instance_type
      vpc_id               = step.create_base_infra.vpc_id
      target_count         = 2
      additional_tags      = step.create_tag1_inputs.tag_map
    }
  }

  step "create_tag2" {
    module = module.random_stringifier
  }

  step "create_tag2_inputs" {
    module     = module.generate_aws_host_tag_vars
    depends_on = [step.create_tag2]

    variables {
      tag_name  = step.create_tag2.string
      tag_value = "test"
    }
  }

  step "create_targets_with_tag2" {
    module     = module.target
    depends_on = [step.create_base_infra]

    variables {
      ami_id               = step.create_base_infra.ami_ids["ubuntu"]["amd64"]
      aws_ssh_keypair_name = var.aws_ssh_keypair_name
      enos_user            = var.enos_user
      instance_type        = var.target_instance_type
      vpc_id               = step.create_base_infra.vpc_id
      target_count         = 1
      additional_tags      = step.create_tag2_inputs.tag_map
    }
  }

  step "create_test_id" {
    module = module.random_stringifier
    variables {
      length = 5
    }
  }

  step "iam_setup" {
    module = module.iam_setup
    depends_on = [
      step.create_base_infra,
      step.create_test_id
    ]

    variables {
      test_id    = step.create_test_id.string
      test_email = var.test_email
    }
  }

  step "run_e2e_test" {
    module = module.test_e2e
    depends_on = [
      step.create_boundary_cluster,
      step.create_targets_with_tag1,
      step.create_targets_with_tag2,
      step.iam_setup
    ]

    variables {
      test_package             = "github.com/hashicorp/boundary/testing/internal/e2e/tests/aws"
      alb_boundary_api_addr    = step.create_boundary_cluster.alb_boundary_api_addr
      auth_method_id           = step.create_boundary_cluster.auth_method_id
      auth_login_name          = step.create_boundary_cluster.auth_login_name
      auth_password            = step.create_boundary_cluster.auth_password
      local_boundary_dir       = local.local_boundary_dir
      aws_ssh_private_key_path = local.aws_ssh_private_key_path
      target_user              = "ubuntu"
      aws_access_key_id        = step.iam_setup.access_key_id
      aws_secret_access_key    = step.iam_setup.secret_access_key
      aws_host_set_filter1     = step.create_tag1_inputs.tag_string
      aws_host_set_ips1        = step.create_targets_with_tag1.target_ips
      aws_host_set_filter2     = step.create_tag2_inputs.tag_string
      aws_host_set_ips2        = step.create_targets_with_tag2.target_ips
    }
  }

  output "test_results" {
    value = step.run_e2e_test.test_results
  }
}
