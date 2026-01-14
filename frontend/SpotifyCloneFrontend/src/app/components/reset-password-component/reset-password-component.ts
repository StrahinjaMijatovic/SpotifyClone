import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-reset-password',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  template: `
    <div class="min-vh-100 d-flex align-items-center bg-black">
      <div class="container">
        <div class="row justify-content-center">
          <div class="col-12 col-md-6 col-lg-4">
            <div class="card bg-dark border-secondary shadow-lg">
              <div class="card-body p-4 p-md-5">
                <div class="text-center mb-4">
                  <h3 class="fw-bold text-light">Nova Lozinka</h3>
                  <p class="text-secondary small">Postavi novu lozinku za svoj nalog.</p>
                </div>

                <div *ngIf="successMessage" class="alert alert-success">
                  {{ successMessage }}
                  <div class="mt-2 text-center">
                    <a routerLink="/login" class="btn btn-sm btn-outline-success">Idi na Login</a>
                  </div>
                </div>
                <div *ngIf="errorMessage" class="alert alert-danger">
                  {{ errorMessage }}
                </div>

                <form (ngSubmit)="onSubmit()" *ngIf="!successMessage">
                   <div *ngIf="!token" class="alert alert-warning">
                      Link za resetovanje je neispravan ili nedostaje.
                   </div>
                  
                  <div class="mb-3" *ngIf="token">
                    <label class="form-label text-secondary">Nova Lozinka</label>
                    <input type="password" class="form-control" [(ngModel)]="newPassword" name="newPassword" required placeholder="Nova lozinka">
                  </div>
                  
                  <button class="btn btn-primary w-100 py-2" type="submit" [disabled]="loading || !token">
                    <span *ngIf="loading" class="spinner-border spinner-border-sm me-2"></span>
                    Promeni Lozinku
                  </button>
                </form>

              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: []
})
export class ResetPasswordComponent implements OnInit {
  token = '';
  newPassword = '';
  loading = false;
  successMessage: string | null = null;
  errorMessage: string | null = null;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private auth: AuthService
  ) { }

  ngOnInit(): void {
    this.route.queryParams.subscribe(params => {
      this.token = params['token'] || '';
      if (!this.token) {
        this.errorMessage = 'Nedostaje token za resetovanje.';
      }
    });
  }

  onSubmit(): void {
    if (!this.token || !this.newPassword) return;

    this.loading = true;
    this.errorMessage = null;

    this.auth.confirmPasswordReset({ token: this.token, new_password: this.newPassword }).subscribe({
      next: () => {
        this.successMessage = 'Lozinka uspešno promenjena! Možeš se ulogovati.';
        this.loading = false;
      },
      error: (err) => {
        this.errorMessage = err?.error?.error || 'Greška pri promeni lozinke. Token je možda istekao.';
        this.loading = false;
      }
    });
  }
}
