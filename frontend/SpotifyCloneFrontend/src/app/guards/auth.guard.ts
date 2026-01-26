import { inject } from '@angular/core';
import { Router, CanActivateFn } from '@angular/router';
import { AuthService } from '../services/auth.service';

export const authGuard: CanActivateFn = (route, state) => {
    const router = inject(Router);
    const authService = inject(AuthService);

    if (authService.isLoggedIn()) {
        return true;
    }

    // Redirect to login page with return url
    return router.createUrlTree(['/login'], { queryParams: { returnUrl: state.url } });
};

export const adminGuard: CanActivateFn = (route, state) => {
    const router = inject(Router);
    const authService = inject(AuthService);

    if (authService.isLoggedIn() && authService.isAdmin()) {
        return true;
    }

    // If logged in but not admin, redirect to home
    if (authService.isLoggedIn()) {
        return router.createUrlTree(['/home']);
    }

    // Not logged in
    return router.createUrlTree(['/login'], { queryParams: { returnUrl: state.url } });
};
