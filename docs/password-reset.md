# Resetting Administrator Password 🔐

This guide details how to recover or reset an administrator account password on **Kodepreneur Panel** if you have forgotten your login credentials.

---

## 🚀 Quick Reset via CLI (Recommended)

Kodepreneur Panel includes a built-in Artisan CLI command to reset any existing administrator password or provision a new Super Administrator account directly from the terminal.

### 1. Interactive Reset Mode

Connect to your server via SSH, navigate to the panel directory, and execute:

```bash
# Production server installation
cd /var/www/kodepreneur-panel
sudo -u www-data php artisan panel:reset-admin
```

*(If running in local development, simply run `php artisan panel:reset-admin` from your `panel/` directory).*

#### How it works:
1. The command displays a list of all existing panel administrator accounts.
2. Prompts you to enter the email address you wish to reset.
3. Prompts for a new password (press **Enter** to automatically generate a secure 16-character random password).
4. Hashes and updates the password in the database immediately.

---

### 2. Direct One-Line Reset

You can also pass the email and new password directly via flags (useful for automated scripts or rapid recovery):

```bash
cd /var/www/kodepreneur-panel
sudo -u www-data php artisan panel:reset-admin admin@kodepreneur.com --password="YourNewSecurePassword123"
```

> **Note:** If the specified email does not exist in the database, the command will offer to create a new user with the `super-admin` role.

---

## 🛠️ Alternative Method 1: Using Laravel Tinker

If you prefer to directly interact with the Eloquent ORM:

1. Open Laravel Tinker:
   ```bash
   cd /var/www/kodepreneur-panel
   sudo -u www-data php artisan tinker
   ```

2. Retrieve the user and set a new hashed password:
   ```php
   // Find the target user by email
   $user = \App\Models\User::where('email', 'admin@kodepreneur.com')->first();

   // Or grab the first administrator account
   // $user = \App\Models\User::first();

   // Update password
   $user->password = \Illuminate\Support\Facades\Hash::make('YourNewSecurePassword123');
   $user->save();
   ```

3. Exit Tinker:
   ```php
   exit
   ```

---

## 🔄 Alternative Method 2: Re-run Default Admin Seeder

If you want to reset the default admin account (`admin@kodepreneur.com`) back to the initial default password (`password`):

```bash
cd /var/www/kodepreneur-panel
sudo -u www-data php artisan db:seed --class=AdminUserSeeder
```

Credentials after running the seeder:
- **Email**: `admin@kodepreneur.com`
- **Password**: `password`

*(Remember to change the default password after logging in!)*

---

## ⚠️ Troubleshooting Common Issues

### Permission Denied on `database.sqlite`
If you encounter a `SQLSTATE[HY000]: General error: 8 attempt to write a readonly database` or file permission error:
Ensure the database file and directory are owned by the web server user:
```bash
sudo chown -R www-data:www-data /var/www/kodepreneur-panel/database
sudo chmod -R 775 /var/www/kodepreneur-panel/database
```

### Rate Limiting / Session Lock
If you failed multiple login attempts before resetting your password, clear Laravel caches to reset temporary throttling:
```bash
sudo -u www-data php artisan cache:clear
sudo -u www-data php artisan optimize:clear
```
