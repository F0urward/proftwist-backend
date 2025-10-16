# How to deploy

```sh
sudo cp deploy/proftwist-back.service /etc/systemd/system/proftwist-back.service
sudo systemctl daemon-reload
sudo systemctl enable --now proftwist-back
```
