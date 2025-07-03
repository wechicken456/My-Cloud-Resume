import { api } from '../api/api';
import { getRecaptchaToken } from '../utils/recaptcha';

const RECAPTCHA_SITE_KEY = 'YOUR_RECAPTCHA_SITE_KEY';

export function setupContactForm(container: HTMLElement) {
    const contactDiv = document.createElement('div');
    contactDiv.id = 'contact-form-board';
    contactDiv.innerHTML = `
        <div id="contact-toggle-btn" aria-label="Contact Me">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
                <polyline points="22,6 12,13 2,6"/>
            </svg>
            <span>Contact</span>
        </div>
        <div id="contact-form-container" style="display: none;">
            <div id="contact-form-header">
                <h3>Contact Me</h3>
                <button id="contact-close-btn" aria-label="Close">&times;</button>
            </div>
            <form id="contact-form">
                <input name="name" placeholder="Your Name" required>
                <input name="email" type="email" placeholder="Your Email" required>
                <textarea name="message" placeholder="Your Message" required rows="3"></textarea>
                <div id="contact-form-actions">
                    <button type="submit">Send Message</button>
                    <span id="contact-status"></span>
                </div>
            </form>
        </div>
    `;
    container.appendChild(contactDiv);

    const toggleBtn = contactDiv.querySelector('#contact-toggle-btn') as HTMLElement;
    const formContainer = contactDiv.querySelector('#contact-form-container') as HTMLElement;
    const closeBtn = contactDiv.querySelector('#contact-close-btn') as HTMLButtonElement;
    const form = contactDiv.querySelector('#contact-form') as HTMLFormElement;
    const status = contactDiv.querySelector('#contact-status') as HTMLElement;

    let isExpanded = false;

    function toggleForm() {
        isExpanded = !isExpanded;
        if (isExpanded) {
            formContainer.style.display = 'block';
            toggleBtn.style.display = 'none';
            contactDiv.classList.add('expanded');
        } else {
            formContainer.style.display = 'none';
            toggleBtn.style.display = 'flex';
            contactDiv.classList.remove('expanded');
            status.textContent = '';
        }
    }

    toggleBtn.onclick = toggleForm;
    closeBtn.onclick = toggleForm;

    form.onsubmit = async (e) => {
        e.preventDefault();
        status.textContent = 'Sending...';
        const formData = new FormData(form);
        const data = {
            name: formData.get('name') as string,
            email: formData.get('email') as string,
            message: formData.get('message') as string,
            recaptcha: await getRecaptchaToken(RECAPTCHA_SITE_KEY)
        };
        const res = await api.sendContact(data);
        if (res.ok) {
            status.textContent = 'Sent!';
            api.sendNotification('contact', { name: data.name, email: data.email });
            form.reset();
            setTimeout(() => {
                toggleForm();
            }, 2000);
        } else {
            status.textContent = 'Error sending message.';
        }
    };
}