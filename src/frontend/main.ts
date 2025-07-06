import './style.css';
import { resumeData } from './data';
import { likeCounter } from './components/Likes';
import { setupContactForm } from './components/Contact';
import { visitorCounter } from './components/Visitor';
import { api } from './api/api';


document.addEventListener('DOMContentLoaded', async () => {
    const app = document.getElementById('app') as HTMLElement;
    visitorCounter.setupVisitorCounter(app);
    likeCounter.setupLikeCounter(app);

    // Header
    const header = document.createElement('header');
    header.className = 'header';
    header.innerHTML = `
        <h1>Tin Vuong</h1>
        <div class="contact-info">
            <a href="mailto:${resumeData.contact.email}">${resumeData.contact.email}</a>
            <span>${resumeData.contact.phone}</span>
            <span>${resumeData.contact.location}</span>
            <a href="${resumeData.contact.github}" target="_blank">${resumeData.contact.github}</a>
        </div>
    `;
    app.appendChild(header);

    // Education
    const educationSection = document.createElement('section');
    educationSection.className = 'section';
    educationSection.innerHTML = `
        <h2>Education</h2>
        <div class="entry">
            <h3>${resumeData.education.institution}</h3>
            <span class="date">${resumeData.education.duration}</span>
            <p>${resumeData.education.degree}</p>
            <p style="margin-bottom: 5px;">${resumeData.education.location}</p>
            <ul>
                ${resumeData.education.details.map(detail => `<li>${detail}</li>`).join('')}
            </ul>
        </div>
    `;
    app.appendChild(educationSection);

    // Experience
    const experienceSection = document.createElement('section');
    experienceSection.className = 'section';
    experienceSection.innerHTML = `
        <h2>Work Experience</h2>
        ${resumeData.experience.map(exp => `
            <div class="entry">
                <h3>${exp.title}</h3>
                <p>${exp.organization}</p>
                <span class="date">${exp.duration}</span>
                <p style="margin-bottom: 5px; font-weight: bold;">${exp.location}</p>
                <ul>
                    ${exp.details.map(detail => `<li>${detail}</li>`).join('')}
                </ul>
            </div>
        `).join('')}
    `;
    app.appendChild(experienceSection);

    // Projects
    const projectsSection = document.createElement('section');
    projectsSection.className = 'section';
    projectsSection.innerHTML = `
        <h2>Projects</h2>
        ${resumeData.projects.map(project => `
            <div class="entry">
                <h3>${project.name}</h3>
                <span class="date">${project.date}</span>
                <ul>
                    ${project.details.map(detail => `<li>${detail}</li>`).join('')}
                </ul>
            </div>
        `).join('')}
    `;
    app.appendChild(projectsSection);

    // Activities
    const activitiesSection = document.createElement('section');
    activitiesSection.className = 'section';
    activitiesSection.innerHTML = `
        <h2>Community & Activities</h2>
        <div class="entry">
            <h3>${resumeData.activities.organization}</h3>
            <span class="date">${resumeData.activities.duration}</span>
            <p>${resumeData.activities.roles.join(', ')}</p>
            <ul>
                ${resumeData.activities.details.map(detail => `<li>${detail}</li>`).join('')}
            </ul>
        </div>
    `;
    app.appendChild(activitiesSection);

    // Skills
    const skillsSection = document.createElement('section');
    skillsSection.className = 'section';

    // Group skills by proficiency
    const proficientSkills = resumeData.skills.filter(skill => skill.proficiency === 'proficient');
    const experiencedSkills = resumeData.skills.filter(skill => skill.proficiency === 'experienced');
    const familiarSkills = resumeData.skills.filter(skill => skill.proficiency === 'familiar');

    skillsSection.innerHTML = `
        <div class="skills-header">
            <h2>Skills</h2>
            <div class="skills-legend">
                <span class="legend-item proficient"><span class="legend-color"></span> Proficient</span>
                <span class="legend-item experienced"><span class="legend-color"></span> Experienced</span>
                <span class="legend-item familiar"><span class="legend-color"></span> Familiar</span>
            </div>
        </div>
        <div class="skills-grid">
            <div class="skill-row proficient">
                ${proficientSkills.map(skill => `<span class="skill-tag skill-${skill.proficiency}">${skill.name}</span>`).join('')}
            </div>
            <div class="skill-row experienced">
                ${experiencedSkills.map(skill => `<span class="skill-tag skill-${skill.proficiency}">${skill.name}</span>`).join('')}
            </div>
            <div class="skill-row familiar">
                ${familiarSkills.map(skill => `<span class="skill-tag skill-${skill.proficiency}">${skill.name}</span>`).join('')}
            </div>
        </div>
    `;
    app.appendChild(skillsSection);


    // Session handling
    // Check if this user has visited or liked before
    api.getSession().then(sessionStatus => {
        // Update visitor and like counters based on this user's session
        visitorCounter.updateVisitorSessionStatus(sessionStatus.has_visited);
        likeCounter.updateLikeSessionStatus(sessionStatus.has_liked);
    }).catch(error => {
        console.error('Failed to fetch session status:', error);
        // Fallback: assume first-time visitor and not liked
        visitorCounter.updateVisitorSessionStatus(false);
        likeCounter.updateLikeSessionStatus(false);
    });
    setupContactForm(app);
});
