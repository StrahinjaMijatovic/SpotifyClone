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

  constructor(
    private contentService: ContentService,
    private authService: AuthService,
    private router: Router
  ) { }

  ngOnInit(): void {
    this.loadContent();
  }

  goToProfile(): void {
    this.router.navigate(['/profile']);
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
