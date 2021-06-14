import React, { createContext, useReducer } from "react";
import { productReducer } from "../reducers/ProductReducer";

export const ProductContext = createContext();

const ProductContextProvider = (props) => {
	const [productState, dispatch] = useReducer(productReducer, {
		listProducts: {
			showError: false,
			errorMessage: "",
			products: [],
		},
		shoppingCart: {},
		createProduct: {
			imageSelected: false,
			showedImage: "",
			showModal: false,
			showErrorMessage: false,
			errorMessage: "",
		},
		updateProduct: {
			imageSelected: false,
			showedImage: "",
			showModal: false,
			showErrorMessage: false,
			errorMessage: "",
			product: {
				id: "",
				price: "",
				name: "",
				quantity: "",
				imageUrl: "",
			},
		},
	});

	return <ProductContext.Provider value={{ productState, dispatch }}>{props.children}</ProductContext.Provider>;
};

export default ProductContextProvider;
