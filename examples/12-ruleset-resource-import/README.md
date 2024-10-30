## Проверка импорта ресурсов

1. Опишите простой ресурс в Терраформ

```
resource "ngenix_ruleset" "terraformimporttestex1" {
  name = "terraformimporttestex1",
}
```

2. Примените манифест - `terraform apply -var-file="vars.tfvars" -auto-approve`

3. Получите ID созданной ДНС зоны с помощью команды `terraform show`

```
# ngenix_ruleset.terraformimporttestex1
resource "ngenix_ruleset" "terraformimporttestex1" {
    id    = "<ID>"
    ...
    last_updated = "Monday, 27-Aug-24 16:18:20 CST"
}
## ...
```

4. Удалить существующуй Ruleset из Terraform state

```
$ terraform state rm ngenix_ruleset.terraformimporttestex1
Removed ngenix_ruleset.terraformimporttestex1
Successfully removed 1 resource instance(s).
```

5. Проверьте что Terraform state больше не содержит информацию о Ruleset. Должна отобразиться только информация из output

```
terraform show

Outputs:

terraformimporttestex1 = {
    id              = "ID"
    name            = "...."
    enabled         = true
    rules           = [
      ...
    ]
##...
```

6. Выполните импорт существующего Ruleset, передав команде импорта идентификатор существующего Ruleset

```
terraform import ngenix_ruleset.terraformimporttestex1 6184123 -var-file="vars.tfvars"
```