import React, { useState } from "react";
import FieldController from "./FieldController";
import { Button } from "@egovernments/digit-ui-react-components";

const MultiChildFormWrapper = ({ config, control, formData, setValue, getValues, errors, props }) => {
  const [instances, setInstances] = useState([{ id: Date.now() }]);

  const addInstance = () => {
    setInstances((prev) => [...prev, { id: Date.now() }]);
  };

  const removeInstance = (id) => {
    setInstances((prev) => prev.filter((inst) => inst.id !== id));
  };

  return (
    <div>
      {instances.map((inst, idx) => (
        <div
          key={inst.id}
          style={{
            position: "relative",
          }}
        >
          {/* Cross (X) icon */}
          {instances.length > 1 && (
            <div
              onClick={() => removeInstance(inst.id)}
              style={{
                position: "absolute",
                top: "10px",
                right: "14px",
                cursor: "pointer",
                fontSize: "18px",
                fontWeight: "bold",
                color: "#555",
              }}
              title="Remove"
            >
              ✕
            </div>
          )}

          {config.body.map((field, i) => {
            const fieldName = `${config.name}.${idx}.${field.populators.name}`;
            return (
              <div style={{marginBottom:"0.75rem"}}>
              <FieldController
                key={`${fieldName}_${i}`}
                type={field.type}
                populators={{ ...field.populators, name: fieldName }}
                isMandatory={field.populators?.required}
                disable={field.populators?.disable}
                component={field.component}
                config={field}
                control={control}
                props={props}
                formData={formData}
                errors={errors}
              />
              </div>
            );
          })}
        </div>
      ))}

      {/* Add button aligned right */}
      {/* <div style={{ marginTop: "20px", display: "flex", justifyContent: "flex-end" }}>
        <Button type="button" variation="secondary" label={`${config?.prefix}_ADD`} onButtonClick={addInstance} />
      </div> */}
    </div>
  );
};

export default MultiChildFormWrapper;
