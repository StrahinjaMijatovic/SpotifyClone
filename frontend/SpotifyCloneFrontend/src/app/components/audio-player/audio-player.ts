import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Subscription } from 'rxjs';
import { AudioPlayerService, PlayerState } from '../../services/audio-player.service';

@Component({
    selector: 'app-audio-player',
    standalone: true,
    imports: [CommonModule, FormsModule],
    templateUrl: './audio-player.html',
    styleUrls: ['./audio-player.css'],
})
export class AudioPlayerComponent implements OnInit, OnDestroy {
    playerState: PlayerState | null = null;
    showQueue = false;
    isLiked = false;
    private subscription?: Subscription;
    private likedSongs: Set<string> = new Set();

    constructor(public audioService: AudioPlayerService) {
        this.loadLikedSongs();
    }

    ngOnInit(): void {
        this.subscription = this.audioService.state$.subscribe(state => {
            this.playerState = state;
            this.updateLikedStatus();
        });
    }

    ngOnDestroy(): void {
        this.subscription?.unsubscribe();
    }

    togglePlayPause(): void {
        this.audioService.togglePlayPause();
    }

    skipNext(): void {
        this.audioService.skipNext();
    }

    skipPrevious(): void {
        this.audioService.skipPrevious();
    }

    onProgressClick(event: MouseEvent): void {
        const progressBar = event.currentTarget as HTMLElement;
        const rect = progressBar.getBoundingClientRect();
        const percent = (event.clientX - rect.left) / rect.width;
        const newTime = percent * (this.playerState?.duration || 0);
        this.audioService.seek(newTime);
    }

    onVolumeChange(event: Event): void {
        const input = event.target as HTMLInputElement;
        this.audioService.setVolume(Number(input.value));
    }

    toggleMute(): void {
        this.audioService.toggleMute();
    }

    toggleShuffle(): void {
        this.audioService.toggleShuffle();
    }

    cycleRepeat(): void {
        this.audioService.cycleRepeatMode();
    }

    get progressPercent(): number {
        if (!this.playerState || !this.playerState.duration) return 0;
        return (this.playerState.currentTime / this.playerState.duration) * 100;
    }

    formatTime(seconds: number): string {
        return this.audioService.formatTime(seconds);
    }

    get repeatIcon(): string {
        switch (this.playerState?.repeatMode) {
            case 'repeat-one': return '🔂';
            case 'repeat-all': return '🔁';
            default: return '🔁';
        }
    }

    get repeatActive(): boolean {
        return this.playerState?.repeatMode !== 'off';
    }

    getArtistsText(): string {
        const artists = this.playerState?.currentSong?.artists;
        if (!artists || artists.length === 0) {
            return 'Nepoznat izvođač';
        }
        return artists.join(', ');
    }

    toggleLike(): void {
        const songId = this.playerState?.currentSong?.id;
        if (!songId) return;

        if (this.likedSongs.has(songId)) {
            this.likedSongs.delete(songId);
            this.isLiked = false;
        } else {
            this.likedSongs.add(songId);
            this.isLiked = true;
        }
        this.saveLikedSongs();
    }

    private loadLikedSongs(): void {
        const stored = localStorage.getItem('likedSongs');
        if (stored) {
            try {
                const arr = JSON.parse(stored);
                this.likedSongs = new Set(arr);
            } catch (e) {
                console.error('Error loading liked songs:', e);
            }
        }
    }

    private saveLikedSongs(): void {
        localStorage.setItem('likedSongs', JSON.stringify([...this.likedSongs]));
    }

    private updateLikedStatus(): void {
        const songId = this.playerState?.currentSong?.id;
        this.isLiked = songId ? this.likedSongs.has(songId) : false;
    }

    toggleQueue(): void {
        this.showQueue = !this.showQueue;
    }

    playFromQueue(index: number): void {
        if (!this.playerState) return;
        const song = this.playerState.queue[index];
        if (song) {
            this.audioService.playSong(song, this.playerState.queue, index, true);
        }
    }

    dismissError(): void {
        this.audioService.clearError();
    }
}
