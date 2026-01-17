import { Component, OnInit, AfterViewInit } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { AuthService } from '../../services/auth.service';
import { RecaptchaService } from '../../services/recaptcha.service';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './register-component.html',
  styleUrl: './register-component.css',
})
export class RegisterComponent implements OnInit, AfterViewInit {
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

  // reCAPTCHA
  recaptchaToken = '';
  recaptchaWidgetId: number | null = null;
  recaptchaLoaded = false;

  // Mapiranje backend poruka na srpski
  private errorMessages: { [key: string]: string } = {
    'Invalid email format': 'Neispravan format email adrese.',
    'This password is too common. Please choose a more unique password': 'Ova lozinka je previše česta. Izaberi jedinstvenu lozinku.',
    'Password must contain uppercase, lowercase, number and special character': 'Lozinka mora sadržati veliko slovo, malo slovo, broj i specijalni znak.',
    'Username must be 3-50 alphanumeric characters or underscore': 'Korisničko ime mora imati 3-50 alfanumeričkih karaktera ili donju crtu.',
    'Names must contain only letters and be 2-50 characters': 'Ime i prezime moraju sadržati samo slova i biti dužine 2-50 karaktera.',
    'Input contains invalid characters': 'Unos sadrži nedozvoljene karaktere.',
    'Username already exists': 'Korisničko ime već postoji.',
    'Email already exists': 'Email adresa već postoji.',
    'Invalid request data': 'Neispravni podaci. Proveri unos.',
    'Invalid reCAPTCHA': 'Neispravan reCAPTCHA. Pokušaj ponovo.',
    'reCAPTCHA verification failed': 'reCAPTCHA verifikacija nije uspela. Pokušaj ponovo.',
  };

  constructor(
    private auth: AuthService,
    private router: Router,
    private recaptchaService: RecaptchaService
  ) {}

  ngOnInit(): void {
    this.recaptchaService.load().then(() => {
      this.recaptchaLoaded = true;
    }).catch(err => {
      console.error('Failed to load reCAPTCHA:', err);
    });
  }

  ngAfterViewInit(): void {
    setTimeout(() => {
      this.renderRecaptcha();
    }, 500);
  }

  renderRecaptcha(): void {
    if (!this.recaptchaLoaded) {
      setTimeout(() => this.renderRecaptcha(), 500);
      return;
    }

    const widgetId = this.recaptchaService.render('recaptcha-register', (token: string) => {
      this.recaptchaToken = token;
    });

    if (widgetId !== null) {
      this.recaptchaWidgetId = widgetId;
    }
  }

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

    // Check reCAPTCHA
    if (!this.recaptchaToken) {
      this.errorMessage = 'Molimo potvrdite da niste robot (reCAPTCHA)';
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
        recaptcha_token: this.recaptchaToken,
      })
      .subscribe({
        next: (res) => {
          this.successMessage = res?.message ?? 'Registracija uspešna! Možeš da se prijaviš.';
          setTimeout(() => this.router.navigate(['/login']), 600);
        },
        error: (err) => {
          const msg = err?.error?.error || err?.error?.message;
          this.errorMessage = this.translateError(msg) || 'Registracija nije uspela. Pokušaj ponovo.';
          this.loading = false;
          // Reset reCAPTCHA on error
          this.recaptchaToken = '';
          if (this.recaptchaWidgetId !== null) {
            this.recaptchaService.reset(this.recaptchaWidgetId);
          }
        },
        complete: () => (this.loading = false),
      });
  }

  private translateError(message: string | undefined): string | null {
    if (!message) return null;
    return this.errorMessages[message] || message;
  }
}
