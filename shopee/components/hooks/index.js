import { useHooks } from "@components/providers/store";

const _isEmpty = (data) => {
    return (
        data == null ||
        data === "" ||
        (Array.isArray(data) && data.length === 0) ||
        (data.constructor === Object && Object.keys(data).length === 0)
    );
};

const enhanceHook = (swrRes) => {
    const { data, error } = swrRes;
    const hasInitialResponse = !!(data || error);
    const isEmpty = hasInitialResponse && _isEmpty(data);

    return {
        ...swrRes,
        isEmpty,
        hasInitialResponse,
    };
};

export const useOrderHistory = (...args) => {
    const swrRes = enhanceHook(useHooks((hooks) => hooks.useOrderHistory)(...args));
    return {
        data: swrRes,
    };
};

export const usePaypal = (...args) => {
    const swrRes = enhanceHook(useHooks((hooks) => hooks.usePaypal)(...args));
    return {
        data: swrRes,
    };
};

export const useSummary = (...args) => {
    const swrRes = enhanceHook(useHooks((hooks) => hooks.useSummary)(...args));
    return {
        data: swrRes,
    };
};

export const useOrders = () => {
    const swrRes = enhanceHook(useHooks((hooks) => hooks.useOrders)());
    return {
        data: swrRes,
    };
};

export const useOrder = (...args) => {
    const swrRes = enhanceHook(useHooks((hooks) => hooks.useOrder)(...args));
    return {
        data: swrRes,
    };
};


export const useProducts = () => {
    const swrRes = enhanceHook(useHooks((hooks) => hooks.useProducts)());
    return {
        data: swrRes,
    };
};

export const useAdminOrders = () => {
    const swrRes = enhanceHook(useHooks((hooks) => hooks.useAdminOrders)());
    return {
        data: swrRes,
    };
};

export const useAdminProducts = () => {
    const swrRes = enhanceHook(useHooks((hooks) => hooks.useAdminProducts)());
    return {
        data: swrRes,
    };
};

export const useAdminProduct = (...args) => {
    const swrRes = enhanceHook(useHooks((hooks) => hooks.useAdminProduct)(...args));
    return {
        data: swrRes,
    };
};

export const useAdminUsers = () => {
    const swrRes = enhanceHook(useHooks((hooks) => hooks.useAdminUsers)());
    return {
        data: swrRes,
    };
};
