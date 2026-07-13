# auto-deploy

```bash
curl -fsSL https://raw.githubusercontent.com/klovbin/auto-deploy/main/install.sh | bash
```

Скрипт сам поставит Go и Git (если нет), скачает репо в `~/auto-deploy`, соберёт и запустит CLI.

Пункты меню: **1** репозиторий → **2** деплой-ключи → **3** клон → **4** CI/CD → **5** ключ `SSH_PRIVATE_KEY` для GitLab.
