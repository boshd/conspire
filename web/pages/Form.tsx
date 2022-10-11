import { FormEvent, useEffect } from "react";
import { useState } from "react";

export const useInput = (initialValue: string) => {
  const [value, setValue] = useState(initialValue);

  return {
    value,
    setValue,
    reset: () => setValue(""),
    bind: {
      value,
      onChange: (event: any) => {
        setValue(event.target.value);
      }
    }
  };
}

export const generateTheory = async (topic: string) => {
  const response = await fetch("http://localhost:8080/theories", {
    method: 'POST',
    body: JSON.stringify({
      topic: topic,
      length: 35,
    })
  })

  return response.text()
}

export function Form(_: any) {
  const { value:topic, bind:bindTopic, reset:resetTopic } = useInput('');

  const [theory, setTheory] = useState('')

  useEffect(() => {

  }, [theory])

  const handleSubmit = async (evt: FormEvent) => {
    evt.preventDefault();
    const response: string = await generateTheory(topic)
    console.log(response)
    setTheory(response)
    resetTopic();
  }
  return (
    <form onSubmit={handleSubmit}>
      <label>
        Topic:
        <input type="text" {...bindTopic} />
      </label>
      <br />
      <br />
      <label>
        Prediciton: {theory}
      </label>
      <br />
      <br />
      <input type="submit" value="Submit" />
    </form>
  );
}

export default Form;