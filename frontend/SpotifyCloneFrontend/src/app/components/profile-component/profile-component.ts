import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="container py-5">
      <div class="row justify-content-center">
        <div class="col-12 col-md-8">
          <div class="d-flex justify-content-between align-items-center mb-4">
            <h2 class="text-light fw-bold">Podešavanja Profila</h2>
            <button class="btn btn-outline-light" (click)="goHome()">
              &larr; Nazad
            </button>
          </div>

          <!-- Profile Details -->
          <div class="card bg-dark border-secondary mb-4">
            <div class="card-body">
              <h5 class="card-title text-light mb-3">Lični Podaci</h5>

              <div *ngIf="profileMessage" class="alert alert-success">
                {{ profileMessage }}
              </div>
              <div *ngIf="profileError" class="alert alert-danger">
                {{ profileError }}
              </div>

              <form (ngSubmit)="onUpdateProfile()" class="row g-3" *ngIf="userProfile">
                <div class="col-md-6">
                  <label class="form-label text-secondary">Ime</label>
                  <input type="text" class="form-control bg-black text-light border-secondary" [(ngModel)]="profileForm.first_name" name="first_name">
                </div>
                <div class="col-md-6">
                  <label class="form-label text-secondary">Prezime</label>
                  <input type="text" class="form-control bg-black text-light border-secondary" [(ngModel)]="profileForm.last_name" name="last_name">
                </div>
                <div class="col-12 mt-4">
                  <button type="submit" class="btn btn-success px-4">Sačuvaj Promene</button>
                </div>
              </form>
            </div>
          </div>

          <!-- Change Password -->
          <div class="card bg-dark border-secondary mb-4">
            <div class="card-body">
              <h5 class="card-title text-light mb-3">Promena Lozinke</h5>
              
              <div *ngIf="passwordMessage" class="alert alert-success">
                {{ passwordMessage }}
              </div>
              <div *ngIf="passwordError" class="alert alert-danger">
                {{ passwordError }}
              </div>

              <form (ngSubmit)="onChangePassword()" class="row g-3">
                <div class="col-12">
                  <label class="form-label text-secondary">Trenutna Lozinka</label>
                  <input type="password" class="form-control bg-black text-light border-secondary" [(ngModel)]="passwordForm.current_password" name="current_password" required>
                </div>
                <div class="col-12">
                  <label class="form-label text-secondary">Nova Lozinka</label>
                  <input type="password" class="form-control bg-black text-light border-secondary" [(ngModel)]="passwordForm.new_password" name="new_password" required>
                </div>
                <!-- Future improvement: confirm new password field -->
                <div class="col-12 mt-4">
                  <button type="submit" class="btn btn-primary px-4">Promeni Lozinku</button>
                </div>
              </form>
            </div>
          </div>

          <!-- Danger Zone -->
          <div class="card border-danger bg-transparent">
            <div class="card-body">
              <h5 class="card-title text-danger mb-3">Opasna Zona</h5>
              <p class="text-secondary small">Ova akcija je nepovratna. Svi tvoji podaci biće trajno obrisani.</p>
              <button class="btn btn-danger" (click)="onDeleteAccount()">
                Obriši Moj Nalog
              </button>
            </div>
          </div>

        </div>
      </div>
    </div>
  `,
  styles: []
})
export class ProfileComponent implements OnInit {
  userProfile: any = null;

  profileForm = { first_name: '', last_name: '' };
  profileMessage: string | null = null;
  profileError: string | null = null;

  passwordForm = { current_password: '', new_password: '' };
  passwordMessage: string | null = null;
  passwordError: string | null = null;

  constructor(private auth: AuthService, private router: Router) { }

  ngOnInit(): void {
    this.loadProfile();
  }

  loadProfile(): void {
    this.auth.getProfile().subscribe({
      next: (user) => {
        this.userProfile = user;
        this.profileForm.first_name = user.first_name;
        this.profileForm.last_name = user.last_name;
      },
      error: () => this.profileError = 'Greška pri učitavanju profila.'
    });
  }

  onUpdateProfile(): void {
    this.profileMessage = null;
    this.profileError = null;

    this.auth.updateProfile(this.profileForm).subscribe({
      next: (res) => {
        this.profileMessage = 'Podaci uspešno ažurirani.';
        setTimeout(() => this.profileMessage = null, 3000);
      },
      error: () => this.profileError = 'Greška pri ažuriranju podataka.'
    });
  }

  onChangePassword(): void {
    this.passwordMessage = null;
    this.passwordError = null;

    if (!this.passwordForm.current_password || !this.passwordForm.new_password) {
      this.passwordError = 'Unesi obe lozinke.';
      return;
    }

    this.auth.changePassword(this.passwordForm).subscribe({
      next: () => {
        this.passwordMessage = 'Lozinka uspešno promenjena.';
        this.passwordForm = { current_password: '', new_password: '' };
      },
      error: (err) => {
        this.passwordError = err?.error?.error || 'Greška pri promeni lozinke.';
      }
    });
  }

  onDeleteAccount(): void {
    if (!confirm('DA LI SI SIGURAN? Ovo briše nalog trajno!')) return;

    this.auth.deleteAccount().subscribe({
      next: () => {
        alert('Nalog obrisan. Doviđenja!');
        this.auth.logout();
        this.router.navigate(['/login']);
      },
      error: () => alert('Desila se greška pri brisanju naloga.')
    });
  }

  goHome(): void {
    this.router.navigate(['/home']);
  }
}
