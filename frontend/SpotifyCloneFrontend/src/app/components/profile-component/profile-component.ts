import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { ContentService } from '../../services/content.service';
import type { Artist, Genre } from '../../models/content.models';

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

          <!-- Subscriptions -->
          <div class="card bg-dark border-secondary mb-4">
            <div class="card-body">
              <h5 class="card-title text-light mb-3">Moje Pretplate</h5>

              <div *ngIf="subscriptionsLoading" class="text-secondary">
                Učitavanje pretplata...
              </div>

              <!-- Subscribed Artists -->
              <div class="mb-4" *ngIf="!subscriptionsLoading">
                <h6 class="text-success mb-2">Umetnici</h6>
                <div *ngIf="subscribedArtists.length === 0" class="text-secondary small">
                  Niste pretplaćeni na nijednog umetnika.
                </div>
                <div class="d-flex flex-wrap gap-2">
                  <div *ngFor="let artist of subscribedArtists"
                       class="badge bg-success bg-opacity-25 text-success border border-success d-flex align-items-center gap-2 p-2">
                    <span style="cursor: pointer;" (click)="goToArtist(artist.id!)">{{ artist.name }}</span>
                    <button class="btn btn-sm btn-outline-danger border-0 p-0 ms-1"
                            (click)="onUnsubscribe(artist.id!, 'artist')"
                            title="Otkaži pretplatu">
                      &times;
                    </button>
                  </div>
                </div>
              </div>

              <!-- Subscribed Genres -->
              <div *ngIf="!subscriptionsLoading">
                <h6 class="text-primary mb-2">Žanrovi</h6>
                <div *ngIf="subscribedGenres.length === 0" class="text-secondary small">
                  Niste pretplaćeni na nijedan žanr.
                </div>
                <div class="d-flex flex-wrap gap-2">
                  <div *ngFor="let genre of subscribedGenres"
                       class="badge bg-primary bg-opacity-25 text-primary border border-primary d-flex align-items-center gap-2 p-2">
                    <span>{{ genre.name }}</span>
                    <button class="btn btn-sm btn-outline-danger border-0 p-0 ms-1"
                            (click)="onUnsubscribe(genre.id!, 'genre')"
                            title="Otkaži pretplatu">
                      &times;
                    </button>
                  </div>
                </div>
              </div>
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

          <div class="text-center mt-4 mb-5">
            <button class="btn btn-secondary px-4" (click)="goHome()">
              &larr; Nazad na Početnu
            </button>
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

  subscribedArtists: Artist[] = [];
  subscribedGenres: Genre[] = [];
  subscriptionsLoading = false;

  constructor(
    private auth: AuthService,
    private router: Router,
    private contentService: ContentService
  ) { }

  ngOnInit(): void {
    this.loadProfile();
    this.loadSubscriptions();
  }

  loadSubscriptions(): void {
    this.subscriptionsLoading = true;
    this.contentService.getSubscriptions().subscribe({
      next: (data) => {
        console.log('Subscriptions loaded:', data);
        this.subscribedArtists = data?.artists ?? [];
        this.subscribedGenres = data?.genres ?? [];
        this.subscriptionsLoading = false;
      },
      error: (err) => {
        console.error('Failed to load subscriptions:', err);
        this.subscriptionsLoading = false;
      }
    });
  }

  onUnsubscribe(targetId: string, type: 'artist' | 'genre'): void {
    this.contentService.unsubscribe(targetId, type).subscribe({
      next: () => {
        if (type === 'artist') {
          this.subscribedArtists = this.subscribedArtists.filter(a => a.id !== targetId);
        } else {
          this.subscribedGenres = this.subscribedGenres.filter(g => g.id !== targetId);
        }
      },
      error: () => {
        alert('Greška pri otkazivanju pretplate.');
      }
    });
  }

  goToArtist(artistId: string): void {
    this.router.navigate(['/artist', artistId]);
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
