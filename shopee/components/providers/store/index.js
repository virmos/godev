import { createContext, useContext, useReducer } from "react";
import { setupHooks } from "./hooks/setupHooks";
import Cookies from "js-cookie";

export const StoreContext = createContext();

const initialState = {
    cart: Cookies.get("cart")
        ? JSON.parse(Cookies.get("cart"))
        : {
              cartItems: [],
              shippingAddress: {},
              paymentMethod: "",
          },
};

function reducer(state, action) {
    switch (action.type) {
        case "CART_ADD_ITEM": {
            const newItem = action.payload;
            const existItem = state.cart.cartItems.find(
                (item) => item.slug === newItem.slug
            );
            const cartItems = existItem
                ? state.cart.cartItems.map((item) =>
                      item.name === existItem.name ? newItem : item
                  )
                : [...state.cart.cartItems, newItem];
            Cookies.set("cart", JSON.stringify({ ...state.cart, cartItems }));
            return {
                ...state,
                cart: {
                    ...state.cart,
                    cartItems,
                },
            };
        }
        case "CART_REMOVE_ITEM": {
            const cartItems = state.cart.cartItems.filter(
                (item) => item.slug !== action.payload.slug
            );
            Cookies.set("cart", JSON.stringify({ ...state.cart, cartItems }));
            return {
                ...state,
                cart: {
                    ...state.cart,
                    cartItems,
                },
            };
        }
        case "CART_RESET":
            Cookies.remove("cart");
            return {
                ...state,
                cart: {
                    cartItems: [],
                    shippingAddress: { location: {} },
                    paymentMethod: "",
                },
            };
        case "CART_CLEAR_ITEMS":
            Cookies.set(
                "cart",
                JSON.stringify({
                    ...state.cart,
                    cartItems: [],
                })
            );
            return {
                ...state,
                cart: {
                    ...state.cart,
                    cartItems: [],
                },
            };

        case "SAVE_SHIPPING_ADDRESS":
            Cookies.set(
                "cart",
                JSON.stringify({
                    ...state.cart,
                    shippingAddress: {
                        ...state.cart.shippingAddress,
                        ...action.payload,
                    },
                })
            );
            return {
                ...state,
                cart: {
                    ...state.cart,
                    shippingAddress: {
                        ...state.cart.shippingAddress,
                        ...action.payload,
                    },
                },
            };
        case "SAVE_PAYMENT_METHOD":
            Cookies.set(
                "cart",
                JSON.stringify({
                    ...state.cart,
                    paymentMethod: action.payload,
                })
            );
            return {
                ...state,
                cart: {
                    ...state.cart,
                    paymentMethod: action.payload,
                },
            };
        default:
            return state;
    }
}

export default function StoreProvider({ children }) {
    const [state, dispatch] = useReducer(reducer, initialState);
    const api = {
        state,
        dispatch,
        hooks: setupHooks(),
    };

    return (
        <StoreContext.Provider value={api}>{children}</StoreContext.Provider>
    );
}

export function useStore() {
    return useContext(StoreContext);
}

export function useHooks(cb) {
    const { hooks } = useStore();
    return cb(hooks);
}