import {  Loader} from "@egovernments/digit-ui-components";
import React from "react";
import { useRouteMatch } from "react-router-dom";
import { default as EmployeeApp } from "./pages/employee";
import PublicServicesCard from "./components/PublicServicesCard";
import { updateCustomConfigs } from "./utils";


// SampleModule component manages the initialization and rendering of the module
export const PublicServicesModule = ({ stateCode, userType, tenants }) => {
  // Get the current route path and URL using React Router
  const { path, url } = useRouteMatch();
  
  // Get the currently selected tenant ID from DIGIT's ULB Service
  const tenantId = Digit.ULBService.getCurrentTenantId();
  
  // Define the modules that this component depends on
  const moduleCode = ["sample", "common", "workflow"];
  
  // Get the current language selected in the DIGIT Store
  const language = Digit.StoreData.getCurrentLanguage();
  
  // Fetch module-specific store data
  const { isLoading, data: store } = Digit.Services.useStore({
    stateCode,
    moduleCode,
    language,
  });

  // Display a loader until the data is available
  if (isLoading) {
    return  <Loader page={true} variant={"PageLoader"}/>;
  }

  // Render the EmployeeApp component with required props
  return <EmployeeApp path={path} stateCode={stateCode} userType={userType} tenants={tenants} />;
};

// Register components to be used in DIGIT's Component Registry
const componentsToRegister = {
  PublicServicesModule,
  PublicServicesCard,
};

// Initialize and register module components
export const initPublicServiceComponents = () => {
  // Apply custom hooks overrides
  
  // Update custom configurations
  updateCustomConfigs();
  
  // Register each component with the DIGIT Component Registry
  Object.entries(componentsToRegister).forEach(([key, value]) => {
    Digit.ComponentRegistryService.setComponent(key, value);
  });
};
