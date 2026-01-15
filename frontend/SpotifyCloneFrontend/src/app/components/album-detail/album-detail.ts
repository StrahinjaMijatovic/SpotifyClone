import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { CommonModule } from '@angular/common';

import { ContentService } from '../../services/content.service';
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

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private contentService: ContentService,
    private cdr: ChangeDetectorRef
  ) { }

  ngOnInit(): void {
    const albumId = this.route.snapshot.paramMap.get('id');
    if (albumId) {
      this.loadAlbum(albumId);
      this.loadSongs(albumId);
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
        this.cdr.detectChanges();
      },
      error: () => {
        this.errorMessage = 'Ne mogu da učitam pesme.';
        this.cdr.detectChanges();
      },
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
