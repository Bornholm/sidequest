import React, {
  ChangeEvent,
  FunctionComponent,
  MouseEvent,
  useCallback,
  useEffect,
  useState,
} from "react";
import { useGenerateCharacter } from "../../hooks/useGenerateCharacter";
import { useNavigate, useParams } from "react-router";
import { Character } from "../../types/Character";
import { useSaveRecord, useRecord } from "../../hooks/useRecord";
import { usePocket } from "../../contexts/pocketbase";

export const EditCharacter: FunctionComponent = () => {
  const params = useParams<{ characterId: string }>();
  const navigate = useNavigate();
  const { user } = usePocket();

  const [regeneratedProperty, setRegeneratedProperty] = useState<string | null>(
    null
  );

  const [character, setCharacter] = useState<Character>({
    age: 0,
    name: "",
    objectives: "",
    race: "",
    sex: "",
    story: "",
  });

  const generateCharacter = useGenerateCharacter();
  const saveCharacter = useSaveRecord("characters");
  const characterRecord = useRecord(
    "characters",
    params.characterId && params.characterId !== "new" ? params.characterId : ""
  );

  useEffect(() => {
    if (!characterRecord.data) return;
    setCharacter({
      id: characterRecord.data.id,
      sex: characterRecord.data.sex,
      name: characterRecord.data.name,
      story: characterRecord.data.story,
      objectives: characterRecord.data.objectives,
      race: characterRecord.data.race,
      age: characterRecord.data.age,
    });
  }, [characterRecord.data]);

  const onCharacterChange = useCallback(
    (evt: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
      const property = evt.currentTarget.dataset.characterProperty;
      const dataType = evt.currentTarget.type;
      if (!property) return;
      let value: any = evt.currentTarget.value;
      switch (dataType) {
        case "number":
          value = parseFloat(value);
          break;
      }
      setCharacter((character) => ({ ...character, [property]: value }));
    },
    []
  );

  const onGenerateClick = useCallback(() => {
    generateCharacter.mutate({});
  }, [character]);

  const onSaveClick = useCallback(() => {
    saveCharacter.mutate({ ...character, author: user?.id });
  }, [character, user]);

  useEffect(() => {
    if (!generateCharacter.data) return;
    const character = generateCharacter.data.character as Character;

    if (regeneratedProperty) {
      setCharacter((originalCharacter) => ({
        ...originalCharacter,
        [regeneratedProperty]: character[regeneratedProperty],
      }));
      setRegeneratedProperty(null);
    } else {
      setCharacter((originalCharacter) => ({
        ...originalCharacter,
        ...character,
      }));
    }
  }, [generateCharacter.data]);

  useEffect(() => {
    if (!saveCharacter.data?.id || character.id) return;
    navigate(`/characters/${saveCharacter.data?.id}`);
  }, [character, saveCharacter.data]);

  const onRegenerateClick = useCallback(
    (evt: MouseEvent<HTMLElement>) => {
      const property = evt.currentTarget.dataset.characterProperty;
      if (!property) return;
      const copy = { ...character };
      delete copy[property];
      setRegeneratedProperty(property);
      generateCharacter.mutate({ character: copy });
    },
    [character]
  );

  const isNew = params.characterId === "new";
  const isLoading =
    characterRecord.isPending ||
    saveCharacter.isPending ||
    generateCharacter.isPending;

  return (
    <div>
      <div className="level is-mobile">
        <div className="level-left">
          <h2 className="title level-item">
            {isNew ? "New character" : `Character "${character.name}"`}
          </h2>
        </div>
        <div className="level-left">
          <button
            disabled={isLoading}
            className={`button is-info level-item ${
              isLoading ? "is-loading" : ""
            }`}
            onClick={onGenerateClick}
          >
            <strong>Generate</strong>
          </button>
          <button
            disabled={isLoading}
            className={`button is-primary level-item ${
              isLoading ? "is-loading" : ""
            }`}
            onClick={onSaveClick}
          >
            <strong>Save</strong>
          </button>
        </div>
      </div>
      <div className="grid">
        <div className="cell">
          <div className="field">
            <label className="label">Name</label>
            <div className="control">
              <input
                className="input"
                type="text"
                placeholder="John Doe"
                data-character-property="name"
                disabled={isLoading}
                onChange={onCharacterChange}
                value={character.name}
              />
              <p className="help">
                <a onClick={onRegenerateClick} data-character-property="name">
                  Regenerate
                </a>
              </p>
            </div>
          </div>
          <div className="field">
            <label className="label">Age</label>
            <div className="control">
              <input
                className="input"
                type="number"
                value={character.age}
                data-character-property="age"
                disabled={isLoading}
                onChange={onCharacterChange}
              />
              <p className="help">
                <a onClick={onRegenerateClick} data-character-property="age">
                  Regenerate
                </a>
              </p>
            </div>
          </div>
          <div className="field">
            <label className="label">Sex</label>
            <div className="control">
              <input
                className="input"
                type="text"
                value={character.sex}
                data-character-property="sex"
                disabled={isLoading}
                onChange={onCharacterChange}
              />
              <p className="help">
                <a onClick={onRegenerateClick} data-character-property="sex">
                  Regenerate
                </a>
              </p>
            </div>
          </div>
          <div className="field">
            <label className="label">Race</label>
            <div className="control">
              <input
                className="input"
                type="text"
                value={character.race}
                data-character-property="race"
                disabled={isLoading}
                onChange={onCharacterChange}
              />
              <p className="help">
                <a onClick={onRegenerateClick} data-character-property="race">
                  Regenerate
                </a>
              </p>
            </div>
          </div>
        </div>
        <div className="cell">
          <div className="field">
            <label className="label">Story</label>
            <div className="control">
              <textarea
                className="textarea"
                rows={15}
                value={character.story}
                data-character-property="story"
                disabled={isLoading}
                onChange={onCharacterChange}
              ></textarea>
              <p className="help">
                <a onClick={onRegenerateClick} data-character-property="story">
                  Regenerate
                </a>
              </p>
            </div>
          </div>
          <div className="field">
            <label className="label">Objectives</label>
            <div className="control">
              <textarea
                className="textarea"
                value={character.objectives}
                data-character-property="objectives"
                disabled={isLoading}
                onChange={onCharacterChange}
              ></textarea>
              <p className="help">
                <a
                  onClick={onRegenerateClick}
                  data-character-property="objectives"
                >
                  Regenerate
                </a>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
