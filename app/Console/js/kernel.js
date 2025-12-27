const { listenerCount } = require("supertest/lib/test");

const kernel = new Kernel();
const router = kernel.getRouter();

document.addEventListener('DOMContentLoaded', () => {
    router.init();
});

document.adoptedStyleSheets([
    new CSSStyleSheet()
])