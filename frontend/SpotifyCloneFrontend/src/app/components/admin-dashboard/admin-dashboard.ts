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

  // Artists
  artists: any[] = [];
  artistForm = { name: '', biography: '', genresInput: '' };
  artistMessage: string | null = null;
  artistError: string | null = null;

  // Albums
  albums: any[] = [];
  albumForm = { name: '', date: '', genre: '', artistsInput: '' };
  albumMessage: string | null = null;
  albumError: string | null = null;

  // Songs
  songs: any[] = [];
  songForm = { name: '', duration: 0, genre: '', album: '', artistsInput: '' };
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

    this.loadArtists();
    this.loadAlbums();
    this.loadSongs();
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

    if (!this.artistForm.name || !this.artistForm.biography || !this.artistForm.genresInput) {
      this.artistError = 'Popuni sva polja!';
      return;
    }

    const genres = this.artistForm.genresInput.split(',').map(g => g.trim());

    this.contentService.createArtist({
      name: this.artistForm.name,
      biography: this.artistForm.biography,
      genres: genres
    }).subscribe({
      next: () => {
        this.artistMessage = 'Umetnik uspešno kreiran!';
        this.artistForm = { name: '', biography: '', genresInput: '' };
        this.loadArtists();
      },
      error: (err) => this.artistError = err?.error?.error || 'Greška pri kreiranju umetnika'
    });
  }

  editArtist(artist: any): void {
    // Pre-populate form for editing
    this.artistForm.name = artist.name;
    this.artistForm.biography = artist.biography;
    this.artistForm.genresInput = artist.genres?.join(', ') || '';

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

    if (!this.albumForm.name || !this.albumForm.date || !this.albumForm.genre || !this.albumForm.artistsInput) {
      this.albumError = 'Popuni sva polja!';
      return;
    }

    const artists = this.albumForm.artistsInput.split(',').map(a => a.trim());

    this.contentService.createAlbum({
      name: this.albumForm.name,
      date: this.albumForm.date,
      genre: this.albumForm.genre,
      artists: artists
    }).subscribe({
      next: () => {
        this.albumMessage = 'Album uspešno kreiran!';
        this.albumForm = { name: '', date: '', genre: '', artistsInput: '' };
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
      !this.songForm.album || !this.songForm.artistsInput) {
      this.songError = 'Popuni sva polja!';
      return;
    }

    const artists = this.songForm.artistsInput.split(',').map(a => a.trim());

    this.contentService.createSong({
      name: this.songForm.name,
      duration: this.songForm.duration,
      genre: this.songForm.genre,
      album: this.songForm.album,
      artists: artists
    }).subscribe({
      next: () => {
        this.songMessage = 'Pesma uspešno kreirana!';
        this.songForm = { name: '', duration: 0, genre: '', album: '', artistsInput: '' };
        this.loadSongs();
      },
      error: (err) => this.songError = err?.error?.error || 'Greška pri kreiranju pesme'
    });
  }
}
