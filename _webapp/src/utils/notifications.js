import store from "@/store"

let defaultTimeout = 5000;

export function addSuccessNotification(text, timeout = defaultTimeout) {
    store.dispatch('app/addNotification', {color: "green", text: text, timeout: timeout}, {root: true});
}

export function addInfoNotification(text, timeout = defaultTimeout) {
    store.dispatch('app/addNotification', {color: "info", text: text, timeout: timeout}, {root: true});
}

export function addPurpleNotification(text, timeout = defaultTimeout) {
    store.dispatch('app/addNotification', {color: "purple", text: text, timeout: timeout}, {root: true});
}

export function addWarningNotification(text, timeout = defaultTimeout) {
    store.dispatch('app/addNotification', {color: "warning", text: text, timeout: timeout}, {root: true});
}

export function addErrorNotificationtext(text, timeout = defaultTimeout) {
    store.dispatch('app/addNotification', {color: "error", text: text, timeout: timeout}, {root: true});
}