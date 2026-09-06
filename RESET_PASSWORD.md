# How to Reset Kodepreneur Panel Password 🔐

If you have forgotten your Kodepreneur Panel administrator credentials, follow the steps below to reset them from your server terminal.

---

## ⚡ Fast CLI Reset (Recommended)

Connect to your server via SSH and execute:

```bash
# In production environment:
cd /var/www/kodepreneur-panel
sudo -u www-data php artisan panel:reset-admin
```

```bash
# In local development:
cd panel
php artisan panel:reset-admin
```

### Options:
1. **Interactive Mode**: Simply run the command without arguments to see a list of accounts and set a new password interactively (or generate one automatically).
2. **One-Line Command**:
   ```bash
   sudo -u www-data php artisan panel:reset-admin admin@kodepreneur.com --password="YourNewPassword123"
   ```

---

## 🛠️ Resetting via Laravel Tinker

```bash
cd /var/www/kodepreneur-panel
sudo -u www-data php artisan tinker
```

Then run inside Tinker:
```php
$user = \App\Models\User::first();
$user->password = \Illuminate\Support\Facades\Hash::make('YourNewPassword123');
$user->save();
exit;
```

---

## 🔄 Reset to Default Admin Seeder

To reset `admin@kodepreneur.com` back to default password (`password`):

```bash
cd /var/www/kodepreneur-panel
sudo -u www-data php artisan db:seed --class=AdminUserSeeder
```

For more details, see [docs/password-reset.md](file:///docs/password-reset.md).
