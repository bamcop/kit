table "todos" {
  schema = schema.main
  column "id" {
    null           = false
    type           = integer
    auto_increment = true
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  column "deleted_at" {
    null = true
    type = datetime
  }
  column "task" {
    null = false
    type = text
  }
  column "completed" {
    null    = false
    type    = bool
    default = false
  }
  primary_key {
    columns = [column.id]
  }
}
table "users" {
  schema = schema.main
  column "id" {
    null           = false
    type           = integer
    auto_increment = true
  }
  column "name" {
    null = false
    type = text
  }
  column "channel" {
    null = false
    type = text
  }
  column "channel_uid" {
    null = false
    type = text
  }
  primary_key {
    columns = [column.id]
  }
  index "user_channel_channel_uid" {
    unique  = true
    columns = [column.channel, column.channel_uid]
  }
}
schema "main" {
}
