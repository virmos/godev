import { handler as createOrderHistoryHook } from "./useOrderHistory";
import { handler as createPaypalHook } from "./usePaypal";
import { handler as createSummaryHook } from "./useSummary";
import { handler as createOrdersHook } from "./useOrders";
import { handler as createOrderHook } from "./useOrder";
import { handler as createProductsHook } from "./useProducts";
import { handler as createAdminOrdersHook } from "./admin/useAdminOrders";
import { handler as createAdminProductsHook } from "./admin/useAdminProducts";
import { handler as createAdminProductHook } from "./admin/useAdminProduct";
import { handler as createAdminUsersHook } from "./admin/useAdminUsers";

export const setupHooks = () => {
    return {
        useOrderHistory: createOrderHistoryHook(),
        usePaypal: createPaypalHook(),
        useSummary: createSummaryHook(),
        useOrders: createOrdersHook(),
        useOrder: createOrderHook(),
        useProducts: createProductsHook(),
        useAdminOrders: createAdminOrdersHook(),
        useAdminProducts: createAdminProductsHook(),
        useAdminProduct: createAdminProductHook(),
        useAdminUsers: createAdminUsersHook(),
    };
};
