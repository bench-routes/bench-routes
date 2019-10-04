interface NotificationsOptions {
    message: string;
    showLoader: boolean;
    criticalLevel: string;
    timeoutInSeconds: number
}

interface LoaderSettings {
    height: number;
}

export class Notifications {
    notificationOpts: NotificationsOptions;
    notificationComponent: HTMLElement;
    loaderPath: string;
    loaderSettings: LoaderSettings;

    constructor() {
        this.notificationComponent = document.getElementById('notification-component');
        this.loaderPath = '../../../assets/img/loader.png';
        this.loaderSettings.height = 20;
    }

    /**
     * Sends a notification to the component along with a self-destruction timer.
     * @param options Notifications options from the interface NotificationsOptions
     */
    send(options: NotificationsOptions): void {
        let notificationElement = document.createElement('span');
        let icon = document.createElement('img'),
            message = document.createElement('span');
        message.innerText = options.message;
        icon.src = this.loaderPath;
        icon.height = this.loaderSettings.height;

        if (options.showLoader) {
            notificationElement.appendChild(icon);
        }
        switch(options.criticalLevel) {
            case 'high':
                message.style.color = 'red';
                break;
            case 'medium':
                message.style.color = 'yellow';
                break;
            case 'low':
                message.style.color = '#444444';
                break;
            default:
                throw new Error('invalid notidication level');
        }
        notificationElement.appendChild(message);
        this.notificationComponent.appendChild(notificationElement);
        setTimeout(() => {
            this.notificationComponent.removeChild(notificationElement);
            notificationElement = null;
        }, 1000 * options.timeoutInSeconds);
    }
}