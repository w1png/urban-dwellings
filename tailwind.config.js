/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.{templ,html,js}"],
  theme: {
    extend: {
      transitionProperty: {
        height: "height",
        width: "width",
      },
      colors: {
        accent: "#6DC8C5",
        delete: "#FF5F5F",
        border: "#8A8A8A",
        hover: "#5d5d5d",
        error: "#FF5F5F",
        "submit-disabled": "#D9D9D9",
        "delete-disabled": "#EE6460",
        "role-user": "#D9D9D9",
        "role-admin": "#FF5F5F",
        "role-artist": "#FFD700",
      },
      dropShadow: {
        glow: [
          "0 0px 20px rgba(255,255, 255, 0.35)",
          "0 0px 65px rgba(255, 255,255, 0.2)",
        ],
      },
    },
  },
  plugins: [],
};
