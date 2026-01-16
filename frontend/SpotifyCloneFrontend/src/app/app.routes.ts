import { Routes } from '@angular/router';
import { HomeComponent } from './components/home-component/home-component';
import { LoginComponent } from './components/login-component/login-component';
import { RegisterComponent } from './components/register-component/register-component';
import { ProfileComponent } from './components/profile-component/profile-component';
import { ForgotPasswordComponent } from './components/forgot-password-component/forgot-password-component';
import { ResetPasswordComponent } from './components/reset-password-component/reset-password-component';
import { MagicLoginComponent } from './components/magic-login-component/magic-login-component';
import { AdminDashboardComponent } from './components/admin-dashboard/admin-dashboard';
import { VerifyEmailComponent } from './components/verify-email/verify-email';
import { ArtistDetailComponent } from './components/artist-detail/artist-detail';
import { AlbumDetailComponent } from './components/album-detail/album-detail';
import { NotificationsComponent } from './components/notifications/notifications';

export const routes: Routes = [
  { path: '', redirectTo: 'home', pathMatch: 'full' },
  { path: 'home', component: HomeComponent },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'profile', component: ProfileComponent },
  { path: 'forgot-password', component: ForgotPasswordComponent },
  { path: 'reset-password', component: ResetPasswordComponent }, // ocekuje token
  { path: 'verify-email', component: VerifyEmailComponent }, // ocekuje token
  { path: 'magic-login', component: MagicLoginComponent }, // ocekuje token
  { path: 'admin', component: AdminDashboardComponent },
  { path: 'artist/:id', component: ArtistDetailComponent },
  { path: 'album/:id', component: AlbumDetailComponent },
  { path: 'notifications', component: NotificationsComponent },
];
