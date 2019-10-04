import { ipcRenderer, remote } from 'electron';
import { Routes } from './routes';
import * as $ from 'jquery';

export default class BenchRoutesUI {
    socket: WebSocket;
    router: Routes;
    currentMainTemplate: HTMLElement;
    navigator: HTMLObjectElement;
    sidebarComponent: HTMLElement;
    sidebarPingNav: HTMLElement;
    sidebarJitterNav: HTMLElement;
    sidebarDashboardNav: HTMLElement;

    constructor() {
        this.currentMainTemplate = document.getElementById('main-component');
        this.navigator = document.createElement('object');
        this.currentMainTemplate.appendChild(this.navigator);
        this.sidebarDashboardNav = document.getElementById('sidebar-dashboard');
        this.sidebarPingNav = document.getElementById('sidebar-ping');
        this.sidebarJitterNav = document.getElementById('sidebar-jitter');
        this.router = new Routes();

        // keep this at last
        this.configInit();
        this.navigationListeners();
    }

    /**
     * Sets the initial configuration of the application.
     */
    configInit() {
        // load home page as the default page in bench-routes
        this.navigator.data = this.router.routes.homePath;
    }

    /**
     * Assigns the navigation listeners to the sidebar component.
     */
    navigationListeners() {
        this.sidebarDashboardNav.onclick = () => {
            this.navigator.data = this.router.routes.homePath;
        };
        this.sidebarPingNav.onclick = () => {
            this.navigator.data = this.router.routes.pingPath;
        };
        this.sidebarJitterNav.onclick = () => {
            this.navigator.data = this.router.routes.jitterPath;
        }
    }

    sidebarAnimation() {
        $(document).ready(() => {
            $('#sidebar-network').show();
            $('#collapseable-network').hide();
            $('#sidebar-network').click(() => {
                $('#collapseable-network').toggle(500);
            });
        });
    }
}

var instance = new BenchRoutesUI();

(function onLoad() {
    instance.sidebarAnimation();
}());
