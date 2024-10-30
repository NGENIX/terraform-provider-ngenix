## Проверка импорта ресурсов

1. Опишите простой ресурс в Терраформ

```
resource "ngenix_dnszone" "terraformimporttestex1" {
  dns_zone = {
    name = "terraformimporttestex1.ru",
    customer_ref = {
      id = 21046
    }
  }
}
```

2. Примените манифест - `terraform apply -var-file="vars.tfvars" -auto-approve`

3. Получите ID созданной ДНС зоны с помощью команды `terraform show`

```
# ngenix_dnszone.terraformimporttestex1
resource "ngenix_dnszone" "terraformimporttestex1" {
    id    = "<ID>"
    dns_zone = [
        ...
    ]
    last_updated = "Monday, 12-Aug-24 11:18:20 CST"
}
## ...
```

4. Удалить существующую ДНС зону из Terraform state

```
$ terraform state rm ngenix_dnszone.terraformimporttestex1
Removed ngenix_dnszone.terraformimporttestex1
Successfully removed 1 resource instance(s).
```

5. Проверьте что Terraform state больше не содержит информацию о ДНС зоне. Должна отобразиться только информация из output

```
terraform show

Outputs:

terraformimporttestex1 = {
    id              = "ID"
    dns_zone        = [
##...
```

6. Выполните импорт существующей ДНС зоны, передав команде импорта идентификатор существующей ДНС зоны

```
terraform import ngenix_dnszone.terraformimporttestex1 6184123 -var-file="vars.tfvars"
```