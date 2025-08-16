import Alpine from 'alpinejs';
import 'htmx.org';
import collapse from '@alpinejs/collapse';
import '../css/main.css'; // Point to your CSS file
// Register plugins
Alpine.plugin(collapse);
// Make Alpine available globally (required for Pines UI)
window.Alpine = Alpine;
Alpine.start();