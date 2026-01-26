import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
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
  templateUrl: './profile-component.html',
  styleUrls: []
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
    private contentService: ContentService,
    private cdr: ChangeDetectorRef
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
        this.cdr.detectChanges();
      },
      error: (err) => {
        console.error('Failed to load subscriptions:', err);
        this.subscriptionsLoading = false;
        this.cdr.detectChanges();
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
        this.cdr.detectChanges();
      },
      error: () => {
        alert('Greška pri otkazivanju pretplate.');
        this.cdr.detectChanges();
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
        this.cdr.detectChanges();
      },
      error: () => {
        this.profileError = 'Greška pri učitavanju profila.';
        this.cdr.detectChanges();
      }
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
