import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { ContentService } from '../../services/content.service';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-admin-dashboard',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './admin-dashboard.html',
  styleUrl: './admin-dashboard.css'
})
export class AdminDashboardComponent implements OnInit {
  activeTab = 'artists';

  // Genres (shared)
  genres: any[] = [];

  // Artists
  artists: any[] = [];
  artistForm = { name: '', biography: '', selectedGenres: [] as string[] };
  artistMessage: string | null = null;
  artistError: string | null = null;

  // Albums
  albums: any[] = [];
  albumForm = { name: '', date: '', genre: '', selectedArtists: [] as string[] };
  albumMessage: string | null = null;
  albumError: string | null = null;

  // Songs
  songs: any[] = [];
  songForm = { name: '', duration: 0, genre: '', album: '', selectedArtists: [] as string[] };
  songMessage: string | null = null;
  songError: string | null = null;

  constructor(
    private contentService: ContentService,
    private authService: AuthService,
    private router: Router
  ) { }

  ngOnInit(): void {
    // Check if user is admin
    if (!this.authService.isAdmin()) {
      alert('Pristup odbijen! Samo admin korisnici mogu pristupiti ovoj stranici.');
      this.router.navigate(['/home']);
      return;
    }

    this.loadGenres();
    this.loadArtists();
    this.loadAlbums();
    this.loadSongs();
  }

  // Genres
  loadGenres(): void {
    this.contentService.getGenres().subscribe({
      next: (data) => this.genres = data,
      error: () => console.error('Greška pri učitavanju žanrova')
    });
  }

  goHome(): void {
    this.router.navigate(['/home']);
  }

  // Artists
  loadArtists(): void {
    this.contentService.getArtists().subscribe({
      next: (data) => this.artists = data,
      error: () => this.artistError = 'Greška pri učitavanju umetnika'
    });
  }

  createArtist(): void {
    this.artistMessage = null;
    this.artistError = null;

    if (!this.artistForm.name || !this.artistForm.biography || this.artistForm.selectedGenres.length === 0) {
      this.artistError = 'Popuni sva polja i izaberi barem jedan žanr!';
      return;
    }

    this.contentService.createArtist({
      name: this.artistForm.name,
      biography: this.artistForm.biography,
      genres: this.artistForm.selectedGenres
    }).subscribe({
      next: () => {
        this.artistMessage = 'Umetnik uspešno kreiran!';
        this.artistForm = { name: '', biography: '', selectedGenres: [] };
        this.loadArtists();
      },
      error: (err) => this.artistError = err?.error?.error || 'Greška pri kreiranju umetnika'
    });
  }

  editArtist(artist: any): void {
    // Pre-populate form for editing
    this.artistForm.name = artist.name;
    this.artistForm.biography = artist.biography;
    this.artistForm.selectedGenres = artist.genres || [];

    // Scroll to form
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }

  // Albums
  loadAlbums(): void {
    this.contentService.getAlbums().subscribe({
      next: (data) => this.albums = data,
      error: () => this.albumError = 'Greška pri učitavanju albuma'
    });
  }

  createAlbum(): void {
    this.albumMessage = null;
    this.albumError = null;

    if (!this.albumForm.name || !this.albumForm.date || !this.albumForm.genre || this.albumForm.selectedArtists.length === 0) {
      this.albumError = 'Popuni sva polja i izaberi barem jednog umetnika!';
      return;
    }

    this.contentService.createAlbum({
      name: this.albumForm.name,
      date: this.albumForm.date + 'T00:00:00Z',
      genre: this.albumForm.genre,
      artists: this.albumForm.selectedArtists
    }).subscribe({
      next: () => {
        this.albumMessage = 'Album uspešno kreiran!';
        this.albumForm = { name: '', date: '', genre: '', selectedArtists: [] };
        this.loadAlbums();
      },
      error: (err) => this.albumError = err?.error?.error || 'Greška pri kreiranju albuma'
    });
  }

  // Songs
  loadSongs(): void {
    this.contentService.getSongs().subscribe({
      next: (data) => this.songs = data,
      error: () => this.songError = 'Greška pri učitavanju pesama'
    });
  }

  createSong(): void {
    this.songMessage = null;
    this.songError = null;

    if (!this.songForm.name || !this.songForm.duration || !this.songForm.genre ||
      !this.songForm.album || this.songForm.selectedArtists.length === 0) {
      this.songError = 'Popuni sva polja i izaberi barem jednog umetnika!';
      return;
    }

    this.contentService.createSong({
      name: this.songForm.name,
      duration: this.songForm.duration,
      genre: this.songForm.genre,
      album: this.songForm.album,
      artists: this.songForm.selectedArtists
    }).subscribe({
      next: () => {
        this.songMessage = 'Pesma uspešno kreirana!';
        this.songForm = { name: '', duration: 0, genre: '', album: '', selectedArtists: [] };
        this.loadSongs();
      },
      error: (err) => this.songError = err?.error?.error || 'Greška pri kreiranju pesme'
    });
  }
}
