import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { CommonModule } from '@angular/common';

import { ContentService } from '../../services/content.service';
import type { Artist, Album } from '../../models/content.models';

@Component({
  selector: 'app-artist-detail',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './artist-detail.html',
  styleUrls: ['./artist-detail.css'],
})
export class ArtistDetailComponent implements OnInit {
  artist: Artist | null = null;
  albums: Album[] = [];
  loading = true;
  errorMessage: string | null = null;
  isSubscribed = false;
  subscriptionLoading = false;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private contentService: ContentService,
    private cdr: ChangeDetectorRef
  ) { }

  ngOnInit(): void {
    const artistId = this.route.snapshot.paramMap.get('id');
    if (artistId) {
      this.loadArtist(artistId);
      this.loadAlbums(artistId);
      this.checkSubscription(artistId);
    }
  }

  checkSubscription(artistId: string): void {
    this.contentService.checkSubscription(artistId, 'artist').subscribe({
      next: (data) => {
        this.isSubscribed = data.subscribed;
      },
      error: () => {
        this.isSubscribed = false;
      }
    });
  }

  toggleSubscription(): void {
    if (!this.artist?.id) return;

    this.subscriptionLoading = true;
    const artistId = this.artist.id;

    if (this.isSubscribed) {
      this.contentService.unsubscribe(artistId, 'artist').subscribe({
        next: () => {
          this.isSubscribed = false;
          this.subscriptionLoading = false;
        },
        error: () => {
          this.subscriptionLoading = false;
        }
      });
    } else {
      this.contentService.subscribe(artistId, 'artist').subscribe({
        next: () => {
          this.isSubscribed = true;
          this.subscriptionLoading = false;
        },
        error: () => {
          this.subscriptionLoading = false;
        }
      });
    }
  }

  loadArtist(id: string): void {
    this.contentService.getArtist(id).subscribe({
      next: (data) => {
        this.artist = data;
        this.loading = false;
        this.cdr.detectChanges();
      },
      error: (err) => {
        this.errorMessage = 'Ne mogu da učitam umetnika: ' + (err?.error?.error || err?.message || 'Nepoznata greška');
        this.loading = false;
        this.cdr.detectChanges();
      },
    });
  }

  loadAlbums(artistId: string): void {
    this.contentService.getAlbumsByArtist(artistId).subscribe({
      next: (data) => {
        this.albums = data ?? [];
        this.cdr.detectChanges();
      },
      error: () => {
        this.errorMessage = 'Ne mogu da učitam albume.';
        this.cdr.detectChanges();
      },
    });
  }

  goToAlbum(albumId: string): void {
    this.router.navigate(['/album', albumId]);
  }

  goBack(): void {
    this.router.navigate(['/home']);
  }
}
