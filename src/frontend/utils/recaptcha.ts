export function loadRecaptcha(siteKey: string) {
    return new Promise<void>((resolve) => {
        if (document.getElementById('recaptcha-script')) return resolve();
        const script = document.createElement('script');
        script.id = 'recaptcha-script';
        script.src = `https://www.google.com/recaptcha/api.js?render=${siteKey}`;
        script.onload = () => resolve();
        document.head.appendChild(script);
    });
}

export async function getRecaptchaToken(siteKey: string) {
    await loadRecaptcha(siteKey);
    // @ts-ignore
    return grecaptcha.execute(siteKey, { action: 'contact' });
}
