import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { ContentService } from '../../services/content.service';
import { AuthService } from '../../services/auth.service';
import type { Album, Artist, Song } from '../../models/content.models';

@Component({
  selector: 'app-home-component',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './home-component.html',
  styleUrls: ['./home-component.css'],
})
export class HomeComponent implements OnInit {
  searchQuery = '';

  artists: Artist[] = [];
  albums: Album[] = [];
  songs: Song[] = [];

  loading = false;
  errorMessage: string | null = null;

  // Profile state
  userProfile: any = null;
  profileForm = {
    first_name: '',
    last_name: ''
  };
  profileMessage: string | null = null;

  constructor(
    private contentService: ContentService,
    private authService: AuthService, // Injected AuthService
    private router: Router
  ) { }

  ngOnInit(): void {
    this.loadContent();
    this.loadProfile();
  }

  loadProfile(): void {
    this.authService.getProfile().subscribe({
      next: (user) => {
        this.userProfile = user;
        this.profileForm.first_name = user.first_name;
        this.profileForm.last_name = user.last_name;
      },
      error: () => console.error('Failed to load profile')
    });
  }

  onUpdateProfile(): void {
    const { first_name, last_name } = this.profileForm;
    if (!first_name || !last_name) return;

    this.authService.updateProfile({ first_name, last_name }).subscribe({
      next: (res) => {
        this.profileMessage = 'Profil uspešno ažuriran.';
        this.userProfile.first_name = first_name;
        this.userProfile.last_name = last_name;
        setTimeout(() => (this.profileMessage = null), 3000);
      },
      error: () => (this.errorMessage = 'Greška pri ažuriranju profila.')
    });
  }

  onDeleteAccount(): void {
    if (!confirm('Da li si siguran da želiš da obrišeš nalog? Ovo je nepovratno!')) return;

    this.authService.deleteAccount().subscribe({
      next: () => {
        alert('Nalog obrisan.');
        this.logout();
      },
      error: () => alert('Greška pri brisanju naloga.')
    });
  }

  loadContent(): void {
    this.loading = true;
    this.errorMessage = null;

    this.contentService.getArtists().subscribe({
      next: (data) => (this.artists = data ?? []),
      error: () => (this.errorMessage = 'Ne mogu da učitam umetnike.'),
    });

    this.contentService.getAlbums().subscribe({
      next: (data) => (this.albums = data ?? []),
      error: () => (this.errorMessage = 'Ne mogu da učitam albume.'),
    });

    this.contentService.getSongs().subscribe({
      next: (data) => (this.songs = data ?? []),
      error: () => (this.errorMessage = 'Ne mogu da učitam pesme.'),
      complete: () => (this.loading = false),
    });
  }

  search(): void {
    const q = this.searchQuery.trim();
    if (!q) return;

    this.loading = true;
    this.errorMessage = null;

    this.contentService.search(q).subscribe({
      next: (data) => {
        this.artists = data?.artists ?? [];
        this.albums = data?.albums ?? [];
        this.songs = data?.songs ?? [];
      },
      error: () => (this.errorMessage = 'Pretraga nije uspela.'),
      complete: () => (this.loading = false),
    });
  }

  clearSearch(): void {
    this.searchQuery = '';
    this.loadContent();
  }

  logout(): void {
    localStorage.removeItem('token');
    this.router.navigate(['/login']);
  }
}
