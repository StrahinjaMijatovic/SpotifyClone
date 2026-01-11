import { Component } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './login-component.html',
  styleUrl: './login-component.css',
})
export class LoginComponent {
  username = '';
  password = '';

  loading = false;
  errorMessage: string | null = null;
  showPassword = false;

  year = new Date().getFullYear();

  constructor(private auth: AuthService, private router: Router) {}

  onSubmit(): void {
    this.errorMessage = null;

    const u = this.username.trim();
    const p = this.password;

    if (!u || !p) {
      this.errorMessage = 'Unesi username i password.';
      return;
    }

    this.loading = true;

    this.auth.login({ username: u, password: p }).subscribe({
      next: () => this.router.navigate(['/home']),
      error: () => {
        this.errorMessage = 'Login nije uspeo. Proveri kredencijale.';
        this.loading = false;
      },
      complete: () => (this.loading = false),
    });
  }
}
