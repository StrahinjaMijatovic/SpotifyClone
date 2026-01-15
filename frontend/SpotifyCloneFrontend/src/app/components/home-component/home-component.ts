import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { ContentService } from '../../services/content.service';
import { AuthService } from '../../services/auth.service';
import type { Artist } from '../../models/content.models';

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
  loading = false;
  errorMessage: string | null = null;

  get isAdmin(): boolean {
    return this.authService.isAdmin();
  }

  constructor(
    private contentService: ContentService,
    private authService: AuthService,
    private router: Router
  ) { }

  ngOnInit(): void {
    this.loadArtists();
  }

  goToProfile(): void {
    this.router.navigate(['/profile']);
  }

  goToAdmin(): void {
    this.router.navigate(['/admin']);
  }

  goToArtist(artistId: string): void {
    this.router.navigate(['/artist', artistId]);
  }

  loadArtists(): void {
    this.loading = true;
    this.errorMessage = null;

    this.contentService.getArtists().subscribe({
      next: (data) => {
        this.artists = data ?? [];
        this.loading = false;
      },
      error: () => {
        this.errorMessage = 'Ne mogu da učitam umetnike.';
        this.loading = false;
      },
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
        this.loading = false;
      },
      error: () => {
        this.errorMessage = 'Pretraga nije uspela.';
        this.loading = false;
      },
    });
  }

  clearSearch(): void {
    this.searchQuery = '';
    this.loadArtists();
  }

  logout(): void {
    localStorage.removeItem('token');
    this.router.navigate(['/login']);
  }
}
