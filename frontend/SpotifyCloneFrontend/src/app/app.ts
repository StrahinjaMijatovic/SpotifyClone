import { Component, signal } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { AudioPlayerComponent } from './components/audio-player/audio-player';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, AudioPlayerComponent],
  templateUrl: './app.html',
  styleUrls: ['./app.css']
})
export class App {
  protected readonly title = signal('SpotifyCloneFrontend');
}
