import { Injectable } from '@angular/core';

declare global {
  interface Window {
    grecaptcha: any;
    onRecaptchaLoad: () => void;
  }
}

@Injectable({
  providedIn: 'root'
})
export class RecaptchaService {
  private siteKey = '6Le-GU0sAAAAAKulZPg-IiJbh8pRFO7lvFjxW6X_';
  private loaded = false;
  private loadPromise: Promise<void> | null = null;

  /**
   * Load reCAPTCHA script dynamically
   */
  load(): Promise<void> {
    if (this.loaded) {
      return Promise.resolve();
    }

    if (this.loadPromise) {
      return this.loadPromise;
    }

    this.loadPromise = new Promise((resolve, reject) => {
      // Define callback
      window.onRecaptchaLoad = () => {
        this.loaded = true;
        resolve();
      };

      // Create script element
      const script = document.createElement('script');
      script.src = `https://www.google.com/recaptcha/api.js?onload=onRecaptchaLoad&render=explicit`;
      script.async = true;
      script.defer = true;
      script.onerror = () => reject(new Error('Failed to load reCAPTCHA'));

      document.head.appendChild(script);
    });

    return this.loadPromise;
  }

  /**
   * Render reCAPTCHA widget in specified container
   */
  render(containerId: string, callback: (token: string) => void): number | null {
    if (!this.loaded || !window.grecaptcha) {
      console.error('reCAPTCHA not loaded');
      return null;
    }

    return window.grecaptcha.render(containerId, {
      sitekey: this.siteKey,
      callback: callback,
      'expired-callback': () => callback(''),
      'error-callback': () => callback('')
    });
  }

  /**
   * Get current reCAPTCHA response token
   */
  getResponse(widgetId?: number): string {
    if (!window.grecaptcha) return '';
    return widgetId !== undefined
      ? window.grecaptcha.getResponse(widgetId)
      : window.grecaptcha.getResponse();
  }

  /**
   * Reset reCAPTCHA widget
   */
  reset(widgetId?: number): void {
    if (!window.grecaptcha) return;
    if (widgetId !== undefined) {
      window.grecaptcha.reset(widgetId);
    } else {
      window.grecaptcha.reset();
    }
  }

  /**
   * Get site key (for component use)
   */
  getSiteKey(): string {
    return this.siteKey;
  }
}
