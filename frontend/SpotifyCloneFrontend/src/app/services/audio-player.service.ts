import { Injectable, OnDestroy, inject } from '@angular/core';
import { BehaviorSubject, Observable, interval, Subscription } from 'rxjs';
import type { Song } from '../models/content.models';
import { AuthService } from './auth.service';

export type RepeatMode = 'off' | 'repeat-one' | 'repeat-all';

export interface PlayerState {
    currentSong: Song | null;
    isPlaying: boolean;
    isLoading: boolean;
    currentTime: number;
    duration: number;
    volume: number;
    isMuted: boolean;
    repeatMode: RepeatMode;
    isShuffleOn: boolean;
    queue: Song[];
    currentIndex: number;
    error: string | null;
}

@Injectable({ providedIn: 'root' })
export class AudioPlayerService implements OnDestroy {
    private audio: HTMLAudioElement;
    private originalQueue: Song[] = [];
    private progressSubscription?: Subscription;
    private authSubscription?: Subscription;
    private boundEventHandlers: Map<string, EventListener> = new Map();
    private authService = inject(AuthService);
    private isClearing = false; // Flag za ignorisanje error-a pri čišćenju

    private playerState = new BehaviorSubject<PlayerState>({
        currentSong: null,
        isPlaying: false,
        isLoading: false,
        currentTime: 0,
        duration: 0,
        volume: 70,
        isMuted: false,
        repeatMode: 'off',
        isShuffleOn: false,
        queue: [],
        currentIndex: -1,
        error: null,
    });

    public state$: Observable<PlayerState> = this.playerState.asObservable();

    constructor() {
        this.audio = new Audio();
        this.audio.volume = 0.7;
        this.setupAudioListeners();
        this.startProgressTracking();
        this.loadStateFromStorage();
        this.setupKeyboardShortcuts();
        this.setupAuthListener();
    }

    private setupAuthListener(): void {
        this.authSubscription = this.authService.isLoggedIn$.subscribe(isLoggedIn => {
            if (!isLoggedIn) {
                this.stopAndClear();
            }
        });
    }

    /** Zaustavlja player i briše sve podatke - koristi se pri logout-u */
    public stopAndClear(): void {
        this.isClearing = true;
        this.audio.pause();
        this.audio.removeAttribute('src');
        this.audio.load(); // Reset audio element
        this.originalQueue = [];
        this.updateState({
            currentSong: null,
            isPlaying: false,
            isLoading: false,
            currentTime: 0,
            duration: 0,
            queue: [],
            currentIndex: -1,
            error: null,
        });
        // Reset flag nakon kratke pauze
        setTimeout(() => {
            this.isClearing = false;
        }, 100);
    }

    ngOnDestroy(): void {
        this.cleanup();
    }

    private cleanup(): void {
        this.progressSubscription?.unsubscribe();
        this.authSubscription?.unsubscribe();
        this.boundEventHandlers.forEach((handler, event) => {
            this.audio.removeEventListener(event, handler);
        });
        this.boundEventHandlers.clear();
        document.removeEventListener('keydown', this.handleKeyDown);
        this.saveStateToStorage();
    }

    private setupAudioListeners(): void {
        const onLoadedMetadata = () => {
            this.updateState({ duration: this.audio.duration, isLoading: false });
        };
        this.audio.addEventListener('loadedmetadata', onLoadedMetadata);
        this.boundEventHandlers.set('loadedmetadata', onLoadedMetadata);

        const onEnded = () => {
            this.handleSongEnded();
        };
        this.audio.addEventListener('ended', onEnded);
        this.boundEventHandlers.set('ended', onEnded);

        const onError = (e: Event) => {
            // Ignoriši greške pri namjernom čišćenju (logout)
            if (this.isClearing) return;
            console.error('Audio playback error:', e);
            this.updateState({ isPlaying: false, isLoading: false, error: 'Greška pri učitavanju pesme' });
        };
        this.audio.addEventListener('error', onError);
        this.boundEventHandlers.set('error', onError);

        const onCanPlay = () => {
            this.updateState({ isLoading: false, error: null });
        };
        this.audio.addEventListener('canplay', onCanPlay);
        this.boundEventHandlers.set('canplay', onCanPlay);

        const onLoadStart = () => {
            // Ignoriši pri čišćenju
            if (this.isClearing) return;
            this.updateState({ isLoading: true, error: null });
        };
        this.audio.addEventListener('loadstart', onLoadStart);
        this.boundEventHandlers.set('loadstart', onLoadStart);

        const onWaiting = () => {
            this.updateState({ isLoading: true });
        };
        this.audio.addEventListener('waiting', onWaiting);
        this.boundEventHandlers.set('waiting', onWaiting);

        const onPlaying = () => {
            this.updateState({ isLoading: false, isPlaying: true });
        };
        this.audio.addEventListener('playing', onPlaying);
        this.boundEventHandlers.set('playing', onPlaying);
    }

    private startProgressTracking(): void {
        this.progressSubscription = interval(100).subscribe(() => {
            if (!this.audio.paused) {
                this.updateState({ currentTime: this.audio.currentTime });
            }
        });
    }

    private handleKeyDown = (event: KeyboardEvent): void => {
        // Ignorisi ako je focus na input elementu
        const target = event.target as HTMLElement;
        if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable) {
            return;
        }

        const state = this.playerState.value;
        if (!state.currentSong) return;

        switch (event.code) {
            case 'Space':
                event.preventDefault();
                this.togglePlayPause();
                break;
            case 'ArrowLeft':
                event.preventDefault();
                this.seek(Math.max(0, this.audio.currentTime - 5));
                break;
            case 'ArrowRight':
                event.preventDefault();
                this.seek(Math.min(this.audio.duration, this.audio.currentTime + 5));
                break;
            case 'ArrowUp':
                event.preventDefault();
                this.setVolume(Math.min(100, state.volume + 5));
                break;
            case 'ArrowDown':
                event.preventDefault();
                this.setVolume(Math.max(0, state.volume - 5));
                break;
            case 'KeyM':
                this.toggleMute();
                break;
            case 'KeyS':
                this.toggleShuffle();
                break;
            case 'KeyR':
                this.cycleRepeatMode();
                break;
        }
    };

    private setupKeyboardShortcuts(): void {
        document.addEventListener('keydown', this.handleKeyDown);
    }

    private saveStateToStorage(): void {
        const state = this.playerState.value;
        const storageState = {
            volume: state.volume,
            isMuted: state.isMuted,
            repeatMode: state.repeatMode,
            isShuffleOn: state.isShuffleOn,
        };
        localStorage.setItem('audioPlayerState', JSON.stringify(storageState));
    }

    private loadStateFromStorage(): void {
        const stored = localStorage.getItem('audioPlayerState');
        if (stored) {
            try {
                const parsed = JSON.parse(stored);
                this.audio.volume = (parsed.volume || 70) / 100;
                this.updateState({
                    volume: parsed.volume || 70,
                    isMuted: parsed.isMuted || false,
                    repeatMode: parsed.repeatMode || 'off',
                    isShuffleOn: parsed.isShuffleOn || false,
                });
                if (parsed.isMuted) {
                    this.audio.volume = 0;
                }
            } catch (e) {
                console.error('Error loading player state:', e);
            }
        }
    }

    private handleSongEnded(): void {
        const state = this.playerState.value;

        if (state.repeatMode === 'repeat-one') {
            this.audio.currentTime = 0;
            this.audio.play();
        } else {
            this.skipNext();
        }
    }

    private updateState(partial: Partial<PlayerState>): void {
        this.playerState.next({ ...this.playerState.value, ...partial });
    }

    public playSong(song: Song, queue: Song[] = [], startIndex: number = 0, isInternalNavigation: boolean = false): void {
        const streamUrl = this.getStreamUrl(song);

        this.audio.src = streamUrl;
        this.audio.load();

        // Samo resetuj originalQueue ako nije interni skip (next/previous)
        if (!isInternalNavigation) {
            this.originalQueue = [...queue];
        }

        this.updateState({
            currentSong: song,
            queue: isInternalNavigation ? this.playerState.value.queue : (queue.length > 0 ? queue : [song]),
            currentIndex: startIndex >= 0 ? startIndex : 0,
            currentTime: 0,
            error: null,
        });

        this.audio.play()
            .then(() => this.updateState({ isPlaying: true }))
            .catch(err => {
                console.error('Play error:', err);
                this.updateState({ error: 'Nije moguće pustiti pesmu' });
            });

        this.saveStateToStorage();
    }

    public togglePlayPause(): void {
        if (this.audio.paused) {
            this.audio.play()
                .then(() => this.updateState({ isPlaying: true }))
                .catch(err => console.error('Play error:', err));
        } else {
            this.audio.pause();
            this.updateState({ isPlaying: false });
        }
    }

    public pause(): void {
        this.audio.pause();
        this.updateState({ isPlaying: false });
    }

    public stop(): void {
        this.audio.pause();
        this.audio.currentTime = 0;
        this.updateState({
            isPlaying: false,
            currentTime: 0,
            currentSong: null
        });
    }

    public skipNext(): void {
        const state = this.playerState.value;
        if (state.queue.length === 0) return;

        let nextIndex = state.currentIndex + 1;

        if (nextIndex >= state.queue.length) {
            if (state.repeatMode === 'repeat-all') {
                nextIndex = 0;
            } else {
                this.stop();
                return;
            }
        }

        const nextSong = state.queue[nextIndex];
        this.playSong(nextSong, state.queue, nextIndex, true);
    }

    public skipPrevious(): void {
        const state = this.playerState.value;

        // If more than 3 seconds into the song, restart it
        if (this.audio.currentTime > 3) {
            this.audio.currentTime = 0;
            return;
        }

        if (state.queue.length === 0) return;

        let prevIndex = state.currentIndex - 1;

        if (prevIndex < 0) {
            if (state.repeatMode === 'repeat-all') {
                prevIndex = state.queue.length - 1;
            } else {
                this.audio.currentTime = 0;
                return;
            }
        }

        const prevSong = state.queue[prevIndex];
        this.playSong(prevSong, state.queue, prevIndex, true);
    }

    public seek(time: number): void {
        this.audio.currentTime = time;
        this.updateState({ currentTime: time });
    }

    public setVolume(volume: number): void {
        const vol = Math.max(0, Math.min(100, volume));
        this.audio.volume = vol / 100;
        this.updateState({ volume: vol, isMuted: false });
    }

    public toggleMute(): void {
        const state = this.playerState.value;
        if (state.isMuted) {
            this.audio.volume = state.volume / 100;
            this.updateState({ isMuted: false });
        } else {
            this.audio.volume = 0;
            this.updateState({ isMuted: true });
        }
    }

    public toggleShuffle(): void {
        const state = this.playerState.value;
        const newShuffleState = !state.isShuffleOn;

        if (newShuffleState) {
            // Enable shuffle - randomize queue
            const currentSong = state.currentSong;
            const newQueue = this.shuffleArray([...this.originalQueue]);

            // Make sure current song is first
            if (currentSong) {
                const currentIdx = newQueue.findIndex(s => s.id === currentSong.id);
                if (currentIdx > 0) {
                    [newQueue[0], newQueue[currentIdx]] = [newQueue[currentIdx], newQueue[0]];
                }
            }

            this.updateState({
                isShuffleOn: true,
                queue: newQueue,
                currentIndex: 0,
            });
        } else {
            // Disable shuffle - restore original queue
            const currentSong = state.currentSong;
            const originalIndex = this.originalQueue.findIndex(s => s.id === currentSong?.id);

            this.updateState({
                isShuffleOn: false,
                queue: [...this.originalQueue],
                currentIndex: originalIndex >= 0 ? originalIndex : 0,
            });
        }
    }

    public setRepeatMode(mode: RepeatMode): void {
        this.updateState({ repeatMode: mode });
    }

    public cycleRepeatMode(): void {
        const state = this.playerState.value;
        const modes: RepeatMode[] = ['off', 'repeat-all', 'repeat-one'];
        const currentIdx = modes.indexOf(state.repeatMode);
        const nextMode = modes[(currentIdx + 1) % modes.length];
        this.setRepeatMode(nextMode);
    }

    public addToQueue(song: Song): void {
        const state = this.playerState.value;
        const newQueue = [...state.queue, song];
        this.originalQueue.push(song);
        this.updateState({ queue: newQueue });
    }

    public removeFromQueue(index: number): void {
        const state = this.playerState.value;
        const newQueue = state.queue.filter((_, i) => i !== index);

        let newIndex = state.currentIndex;
        if (index < state.currentIndex) {
            newIndex--;
        } else if (index === state.currentIndex) {
            // If removing current song, stop playback
            this.stop();
            newIndex = -1;
        }

        this.updateState({
            queue: newQueue,
            currentIndex: newIndex
        });
    }

    public clearQueue(): void {
        this.stop();
        this.originalQueue = [];
        this.updateState({
            queue: [],
            currentIndex: -1
        });
    }

    private shuffleArray<T>(array: T[]): T[] {
        const shuffled = [...array];
        for (let i = shuffled.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
        }
        return shuffled;
    }

    private getStreamUrl(song: Song): string {
        // If song has audio_url, use it directly
        if (song.audio_url) {
            return song.audio_url;
        }

        // Otherwise, use the streaming endpoint
        const token = localStorage.getItem('token');
        return `/api/v1/songs/${song.id}/stream?token=${token}`;
    }

    public formatTime(seconds: number): string {
        if (!seconds || isNaN(seconds)) return '0:00';
        const mins = Math.floor(seconds / 60);
        const secs = Math.floor(seconds % 60);
        return `${mins}:${secs.toString().padStart(2, '0')}`;
    }

    public clearError(): void {
        this.updateState({ error: null });
    }
}
