import { handler as createDeliverMutate } from "./useDeliver";
import { handler as createPayMutate } from "./usePay";
import { handler as createOrderHistoryHook } from "./useOrderHistory";
import { handler as createPaypalHook } from "./usePaypal";
import { handler as createSummaryHook } from "./useSummary";
import { handler as createOrderHook } from "./useOrder";
import { handler as createOrdersHook } from "./useOrders";
import { handler as createProductHook } from "./useProduct";
import { handler as createProductsHook } from "./useProducts";
import { handler as createAdminOrderHook } from "./admin/useAdminOrder";
import { handler as createAdminOrdersHook } from "./admin/useAdminOrders";
import { handler as createAdminProductHook } from "./admin/useAdminProduct";
import { handler as createAdminProductsHook } from "./admin/useAdminProducts";
import { handler as createAdminUserHook } from "./admin/useAdminUser";
import { handler as createAdminUsersHook } from "./admin/useAdminUsers";

export const setupHooks = () => {
    return {
        useDeliverMutate: createDeliverMutate(),
        usePayMutate: createPayMutate(),
        useOrderHistory: createOrderHistoryHook(),
        usePaypal: createPaypalHook(),
        useSummary: createSummaryHook(),
        useOrder: createOrderHook(),
        useOrders: createOrdersHook(),
        useProduct: createProductHook(),
        useProducts: createProductsHook(),
        useAdminOrder: createAdminOrderHook(),
        useAdminOrders: createAdminOrdersHook(),
        useAdminProduct: createAdminProductHook(),
        useAdminProducts: createAdminProductsHook(),
        useAdminUser: createAdminUserHook(),
        useAdminUsers: createAdminUsersHook(),
    };
};
