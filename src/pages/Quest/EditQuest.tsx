import React, {
  ChangeEvent,
  FunctionComponent,
  MouseEvent,
  useCallback,
  useEffect,
  useState,
} from "react";
import { useGenerateQuest } from "../../hooks/useGenerateQuest";
import { useNavigate, useParams } from "react-router";
import { Quest } from "../../types/Quest";
import { useSaveRecord, useRecord } from "../../hooks/useRecord";
import { usePocket } from "../../contexts/pocketbase";

export const EditQuest: FunctionComponent = () => {
  const params = useParams<{ questId: string }>();
  const navigate = useNavigate();
  const { user } = usePocket();

  const [regeneratedProperty, setRegeneratedProperty] = useState<string | null>(
    null
  );

  const [quest, setQuest] = useState<Quest>({
    characters: "",
    clues: "",
    title: "",
    description: "",
    solution: "",
  });

  const generateQuest = useGenerateQuest();
  const saveQuest = useSaveRecord("quests");
  const questRecord = useRecord(
    "quests",
    params.questId && params.questId !== "new" ? params.questId : ""
  );

  useEffect(() => {
    if (!questRecord.data) return;
    setQuest({
      id: questRecord.data.id,
      characters: questRecord.data.characters,
      clues: questRecord.data.clues,
      title: questRecord.data.title,
      description: questRecord.data.description,
      solution: questRecord.data.solution,
    });
  }, [questRecord.data]);

  const onQuestChange = useCallback(
    (evt: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
      const property = evt.currentTarget.dataset.questProperty;
      const dataType = evt.currentTarget.type;
      if (!property) return;
      let value: any = evt.currentTarget.value;
      switch (dataType) {
        case "number":
          value = parseFloat(value);
          break;
      }
      setQuest((quest) => ({ ...quest, [property]: value }));
    },
    []
  );

  const onGenerateClick = useCallback(() => {
    generateQuest.mutate({});
  }, [quest]);

  const onSaveClick = useCallback(() => {
    saveQuest.mutate({ ...quest, author: user?.id });
  }, [quest, user]);

  useEffect(() => {
    if (!generateQuest.data) return;
    const quest = generateQuest.data.quest as Quest;

    if (regeneratedProperty) {
      setQuest((originalQuest) => ({
        ...originalQuest,
        [regeneratedProperty]: quest[regeneratedProperty],
      }));
      setRegeneratedProperty(null);
    } else {
      setQuest((originalQuest) => ({
        ...originalQuest,
        ...quest,
      }));
    }
  }, [generateQuest.data]);

  useEffect(() => {
    if (!saveQuest.data?.id || quest.id) return;
    navigate(`/quests/${saveQuest.data?.id}`);
  }, [quest, saveQuest.data]);

  const onRegenerateClick = useCallback(
    (evt: MouseEvent<HTMLElement>) => {
      const property = evt.currentTarget.dataset.questProperty;
      if (!property) return;
      const copy = { ...quest };
      delete copy[property];
      setRegeneratedProperty(property);
      generateQuest.mutate({ quest: copy });
    },
    [quest]
  );

  const isNew = params.questId === "new";
  const isLoading =
    questRecord.isPending || saveQuest.isPending || generateQuest.isPending;

  return (
    <div>
      {generateQuest.error ? (
        <div className="message is-danger">
          <div className="message-body">{generateQuest.error.message}</div>
        </div>
      ) : null}
      <div className="level">
        <div className="level-left">
          <h2 className="title level-item">
            {isNew ? "New quest" : `Quest "${quest.title}"`}
          </h2>
        </div>
        <div className="level-left">
          <div className="buttons">
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
      </div>
      <div className="columns">
        <div className="column">
          <div className="field">
            <label className="label">Title</label>
            <div className="control">
              <input
                className="input"
                type="text"
                placeholder="The quest for the Holy Grail"
                data-quest-property="title"
                disabled={isLoading}
                onChange={onQuestChange}
                value={quest.title}
              />
              <p className="help">
                <a onClick={onRegenerateClick} data-quest-property="title">
                  Regenerate
                </a>
              </p>
            </div>
          </div>
          <div className="field">
            <label className="label">Description</label>
            <div className="control">
              <textarea
                className="textarea"
                rows={15}
                value={quest.description}
                data-quest-property="description"
                disabled={isLoading}
                onChange={onQuestChange}
              ></textarea>
              <p className="help">
                <a
                  onClick={onRegenerateClick}
                  data-quest-property="description"
                >
                  Regenerate
                </a>
              </p>
            </div>
          </div>
        </div>
        <div className="column">
          <div className="field">
            <label className="label">Characters</label>
            <div className="control">
              <textarea
                className="textarea"
                value={quest.characters}
                data-quest-property="characters"
                disabled={isLoading}
                onChange={onQuestChange}
              ></textarea>
              <p className="help">
                <a onClick={onRegenerateClick} data-quest-property="characters">
                  Regenerate
                </a>
              </p>
            </div>
          </div>
          <div className="field">
            <label className="label">Clues</label>
            <div className="control">
              <textarea
                className="textarea"
                value={quest.clues}
                data-quest-property="clues"
                disabled={isLoading}
                onChange={onQuestChange}
              ></textarea>
              <p className="help">
                <a onClick={onRegenerateClick} data-quest-property="clues">
                  Regenerate
                </a>
              </p>
            </div>
          </div>
          <div className="field">
            <label className="label">Solution</label>
            <div className="control">
              <textarea
                className="textarea"
                value={quest.solution}
                data-quest-property="solution"
                disabled={isLoading}
                onChange={onQuestChange}
              ></textarea>
              <p className="help">
                <a onClick={onRegenerateClick} data-quest-property="solution">
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
