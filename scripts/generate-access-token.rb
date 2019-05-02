terraform_token = PersonalAccessToken.create(user_id: 1, scopes: [:api, :read_user], name: :terraform)
terraform_token.set_token(:ACCTEST)
terraform_token.save!
