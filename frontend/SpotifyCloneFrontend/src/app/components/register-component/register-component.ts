import { Component } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './register-component.html',
  styleUrl: './register-component.css',
})
export class RegisterComponent {
  username = '';
  email = '';
  firstName = '';
  lastName = '';
  password = '';
  passwordConfirm = '';

  loading = false;
  errorMessage: string | null = null;
  successMessage: string | null = null;

  showPassword = false;
  showPasswordConfirm = false;

  year = new Date().getFullYear();

  constructor(private auth: AuthService, private router: Router) {}

  onSubmit(): void {
    this.errorMessage = null;
    this.successMessage = null;

    const username = this.username.trim();
    const email = this.email.trim();
    const firstName = this.firstName.trim();
    const lastName = this.lastName.trim();

    if (!username || !email || !firstName || !lastName || !this.password || !this.passwordConfirm) {
      this.errorMessage = 'Popuni sva polja.';
      return;
    }

    if (this.password !== this.passwordConfirm) {
      this.errorMessage = 'Password i potvrda password-a se ne poklapaju.';
      return;
    }

    this.loading = true;

    this.auth
      .register({
        username,
        email,
        first_name: firstName,
        last_name: lastName,
        password: this.password,
        password_confirm: this.passwordConfirm,
      })
      .subscribe({
        next: (res) => {
          this.successMessage = res?.message ?? 'Registracija uspešna! Možeš da se prijaviš.';
          // mali delay nije potreban, ali UX je lepši
          setTimeout(() => this.router.navigate(['/login']), 600);
        },
        error: (err) => {
          // ako backend vrati {error:"..."} pokušaj da prikažeš to
          const msg = err?.error?.error || err?.error?.message;
          this.errorMessage = msg || 'Registracija nije uspela. Pokušaj ponovo.';
          this.loading = false;
        },
        complete: () => (this.loading = false),
      });
  }
}
