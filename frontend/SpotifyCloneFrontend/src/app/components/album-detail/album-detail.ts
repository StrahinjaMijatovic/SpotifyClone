import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { CommonModule } from '@angular/common';

import { ContentService } from '../../services/content.service';
import { RatingsService, Rating } from '../../services/ratings.service';
import type { Album, Song } from '../../models/content.models';

@Component({
  selector: 'app-album-detail',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './album-detail.html',
  styleUrls: ['./album-detail.css'],
})
export class AlbumDetailComponent implements OnInit {
  album: Album | null = null;
  songs: Song[] = [];
  loading = true;
  errorMessage: string | null = null;

  // Ratings
  myRatings: Map<string, number> = new Map(); // song_id -> rating
  songAverages: Map<string, number> = new Map(); // song_id -> average
  songCounts: Map<string, number> = new Map(); // song_id -> count
  ratingInProgress: string | null = null; // song_id currently being rated

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private contentService: ContentService,
    private ratingsService: RatingsService,
    private cdr: ChangeDetectorRef
  ) { }

  ngOnInit(): void {
    const albumId = this.route.snapshot.paramMap.get('id');
    if (albumId) {
      this.loadAlbum(albumId);
      this.loadSongs(albumId);
      this.loadMyRatings();
    }
  }

  loadAlbum(id: string): void {
    this.contentService.getAlbum(id).subscribe({
      next: (data) => {
        this.album = data;
        this.loading = false;
        this.cdr.detectChanges();
      },
      error: () => {
        this.errorMessage = 'Ne mogu da učitam album.';
        this.loading = false;
        this.cdr.detectChanges();
      },
    });
  }

  loadSongs(albumId: string): void {
    this.contentService.getSongsByAlbum(albumId).subscribe({
      next: (data) => {
        this.songs = data ?? [];
        // Load ratings for each song
        this.songs.forEach(song => {
          if (song.id) {
            this.loadSongRatings(song.id);
          }
        });
        this.cdr.detectChanges();
      },
      error: () => {
        this.errorMessage = 'Ne mogu da učitam pesme.';
        this.cdr.detectChanges();
      },
    });
  }

  loadMyRatings(): void {
    this.ratingsService.getMyRatings().subscribe({
      next: (ratings) => {
        this.myRatings.clear();
        ratings?.forEach(r => {
          this.myRatings.set(r.song_id, r.rating);
        });
        this.cdr.detectChanges();
      },
      error: () => { }
    });
  }

  loadSongRatings(songId: string): void {
    this.ratingsService.getSongRatings(songId).subscribe({
      next: (data) => {
        this.songAverages.set(songId, data.average);
        this.songCounts.set(songId, data.count);
        this.cdr.detectChanges();
      },
      error: () => { }
    });
  }

  getMyRating(songId: string): number {
    return this.myRatings.get(songId) || 0;
  }

  getSongAverage(songId: string): number {
    return this.songAverages.get(songId) || 0;
  }

  getSongCount(songId: string): number {
    return this.songCounts.get(songId) || 0;
  }

  rateSong(songId: string, rating: number): void {
    this.ratingInProgress = songId;
    this.ratingsService.rateSong(songId, rating).subscribe({
      next: () => {
        this.myRatings.set(songId, rating);
        this.loadSongRatings(songId); // Refresh average
        this.ratingInProgress = null;
        this.cdr.detectChanges();
      },
      error: () => {
        this.ratingInProgress = null;
        this.cdr.detectChanges();
      }
    });
  }

  deleteRating(songId: string): void {
    this.ratingInProgress = songId;
    this.ratingsService.deleteRating(songId).subscribe({
      next: () => {
        this.myRatings.delete(songId);
        this.loadSongRatings(songId); // Refresh average
        this.ratingInProgress = null;
        this.cdr.detectChanges();
      },
      error: () => {
        this.ratingInProgress = null;
        this.cdr.detectChanges();
      }
    });
  }

  formatDuration(seconds: number): string {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  }

  goBack(): void {
    window.history.back();
  }

  goHome(): void {
    this.router.navigate(['/home']);
  }
}
